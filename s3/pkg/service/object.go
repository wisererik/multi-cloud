// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/journeymidnight/yig/helper"
	"github.com/micro/go-micro/metadata"
	"github.com/opensds/multi-cloud/api/pkg/common"
	"github.com/opensds/multi-cloud/api/pkg/utils/constants"
	"github.com/opensds/multi-cloud/s3/error"
	. "github.com/opensds/multi-cloud/s3/error"
	dscommon "github.com/opensds/multi-cloud/s3/pkg/datastore/common"
	"github.com/opensds/multi-cloud/s3/pkg/datastore/driver"
	"github.com/opensds/multi-cloud/s3/pkg/meta/types"
	. "github.com/opensds/multi-cloud/s3/pkg/meta/types"
	meta "github.com/opensds/multi-cloud/s3/pkg/meta/types"
	"github.com/opensds/multi-cloud/s3/pkg/meta/util"
	"github.com/opensds/multi-cloud/s3/pkg/utils"
	pb "github.com/opensds/multi-cloud/s3/proto"
	log "github.com/sirupsen/logrus"
)

var ChunkSize int = 2048

func (s *s3Service) CreateObject(ctx context.Context, in *pb.Object, out *pb.BaseResponse) error {
	log.Infoln("CreateObject is called in s3 service.")

	return nil
}

func (s *s3Service) UpdateObject(ctx context.Context, in *pb.Object, out *pb.BaseResponse) error {
	log.Infoln("PutObject is called in s3 service.")

	return nil
}

type DataStreamRecv interface {
	Recv() (*pb.PutDataStream, error)
}

type StreamReader struct {
	in   DataStreamRecv
	req  *pb.PutDataStream
	curr int
}

func (dr *StreamReader) Read(p []byte) (n int, err error) {
	left := len(p)
	for left > 0 {
		if dr.curr == 0 || (dr.req != nil && dr.curr == len(dr.req.Data)) {
			dr.req, err = dr.in.Recv()
			if err != nil && err != io.EOF {
				log.Errorln("failed to recv data with err:", err)
				return
			}
			if len(dr.req.Data) == 0 {
				log.Errorln("no data left to read.")
				err = io.EOF
				return
			}
			dr.curr = 0
		}

		copyLen := 0
		if len(dr.req.Data)-dr.curr > left {
			copyLen = left
		} else {
			copyLen = len(dr.req.Data) - dr.curr
		}
		log.Traceln("copy len:", copyLen)
		copy(p[n:], dr.req.Data[dr.curr:(dr.curr+copyLen)])
		dr.curr += copyLen
		left -= copyLen
		n += copyLen
	}
	return
}

func (s *s3Service) PutObject(ctx context.Context, in pb.S3_PutObjectStream) error {
	log.Infoln("PutObject is called in s3 service.")

	result := &pb.PutObjectResponse{}
	defer in.SendMsg(result)

	var ok bool
	var tenantId string
	var md map[string]string
	md, ok = metadata.FromContext(ctx)
	if !ok {
		log.Error("get metadata from ctx failed.")
		return ErrInternalError
	}

	if tenantId, ok = md[common.CTX_KEY_TENANT_ID]; !ok {
		log.Error("get tenantid failed.")
		return ErrInternalError
	}

	obj := &pb.Object{}
	err := in.RecvMsg(obj)
	if err != nil {
		log.Errorln("failed to get msg with err:", err)
		return ErrInternalError
	}
	obj.TenantId = tenantId
	log.Infof("metadata of object is:%+v\n", obj)
	log.Infof("*********bucket:%s,key:%s,size:%d\n", obj.BucketName, obj.ObjectKey, obj.Size)
	oldObj, err := s.MetaStorage.GetObject(ctx, obj.BucketName, obj.ObjectKey, "", true)
	if err != nil && oldObj != nil {
		log.Info("got the object with same name(%s, %s, %s)", oldObj.BucketName, oldObj.ObjectKey, oldObj.ObjectId)
		obj.StorageMeta = oldObj.StorageMeta
		obj.ObjectId = oldObj.ObjectId
	}

	bucket, err := s.MetaStorage.GetBucket(ctx, obj.BucketName, true)
	if err != nil {
		log.Errorln("get bucket failed with err:", err)
		return err
	}
	if bucket == nil {
		log.Infoln("bucket is nil")
	}

	switch bucket.Acl.CannedAcl {
	case "public-read-write":
		break
	default:
		if bucket.TenantId != obj.TenantId {
			return s3error.ErrBucketAccessForbidden
		}
	}

	data := &StreamReader{in: in}
	var limitedDataReader io.Reader
	if obj.Size > 0 { // request.ContentLength is -1 if length is unknown
		limitedDataReader = io.LimitReader(data, obj.Size)
	} else {
		limitedDataReader = data
	}

	backendName := bucket.DefaultLocation
	if obj.Location != "" {
		backendName = obj.Location
	}
	backend, err := utils.GetBackend(ctx, s.backendClient, backendName)
	if err != nil {
		log.Errorln("failed to get backend client with err:", err)
		return err
	}

	log.Infoln("bucket location:", obj.Location, " backendtype:", backend.Type,
		" endpoint:", backend.Endpoint)
	bodyMd5 := md["md5Sum"]
	ctx = context.Background()
	ctx = context.WithValue(ctx, dscommon.CONTEXT_KEY_SIZE, obj.Size)
	ctx = context.WithValue(ctx, dscommon.CONTEXT_KEY_MD5, bodyMd5)
	sd, err := driver.CreateStorageDriver(backend.Type, backend)
	if err != nil {
		log.Errorln("failed to create storage. err:", err)
		return nil
	}
	res, err := sd.Put(ctx, limitedDataReader, obj)
	if err != nil {
		log.Errorln("failed to put data. err:", err)
		return err
	}
	oid := res.ObjectId
	bytesWritten := res.Written

	// Should metadata update failed, add `maybeObjectToRecycle` to `RecycleQueue`,
	// so the object in Ceph could be removed asynchronously
	//maybeObjectToRecycle := objectToRecycle{
	//	location: cephCluster.Name,
	//	pool:     poolName,
	//	objectId: oid,
	//}
	if bytesWritten < obj.Size {
		//RecycleQueue <- maybeObjectToRecycle
		log.Warnf("write objects, already written(%d), total size(%d)\n", bytesWritten, obj.Size)
	}
	result.Md5 = res.Etag

	/*if signVerifyReader, ok := data.(*signature.SignVerifyReader); ok {
		credential, err = signVerifyReader.Verify()
		if err != nil {
			RecycleQueue <- maybeObjectToRecycle
			return
		}
	}*/
	// TODO validate bucket policy and fancy ACL
	obj.ObjectId = oid
	obj.LastModified = time.Now().UTC().Unix()
	obj.Etag = res.Etag
	obj.ContentType = md["Content-Type"]
	obj.DeleteMarker = false
	obj.CustomAttributes = md /* TODO: only reserve http header attr*/
	obj.Type = meta.ObjectTypeNormal
	obj.Tier = utils.Tier1 // Currently only support tier1
	obj.StorageMeta = res.Meta

	object := &meta.Object{Object: obj}

	result.LastModified = object.LastModified
	//var nullVerNum uint64
	//nullVerNum, err = yig.checkOldObject(bucketName, objectName, bucket.Versioning)
	//if err != nil {
	//	RecycleQueue <- maybeObjectToRecycle
	//	return
	//}
	//if bucket.Versioning == "Enabled" {
	//	result.VersionId = object.GetVersionId()
	//}
	//// update null version number
	//if bucket.Versioning == "Suspended" {
	//	nullVerNum = uint64(object.LastModifiedTime.UnixNano())
	//}
	//
	//if nullVerNum != 0 {
	//	objMap := &meta.ObjMap{
	//		Name:       objectName,
	//		BucketName: bucketName,
	//	}
	//	err = yig.MetaStorage.PutObject(object, nil, objMap, true)
	//} else {
	//	err = yig.MetaStorage.PutObject(object, nil, nil, true)
	//}

	err = s.MetaStorage.PutObject(ctx, object, nil, nil, true)
	if err != nil {
		log.Errorln("failed to put object meta. err:", err)
		//RecycleQueue <- maybeObjectToRecycle
		return ErrDBError
	}
	if err == nil {
		//b.MetaStorage.Cache.Remove(redis.ObjectTable, obj.OBJECT_CACHE_PREFIX, bucketName+":"+objectKey+":")
		//b.DataCache.Remove(bucketName + ":" + objectKey + ":" + object.GetVersionId())
	}

	return nil
}

func (s *s3Service) GetObjectMeta(ctx context.Context, in *pb.Object, out *pb.GetObjectMetaResult) error {
	var err error
	defer func() {
		out.ErrorCode = GetErrCode(err)
	}()

	object, err := s.MetaStorage.GetObject(ctx, in.BucketName, in.ObjectKey, in.VersionId, true)
	if err != nil {
		log.Errorln("failed to get object info from meta storage. err:", err)
		return nil
	}

	_, _, err = CheckRights(ctx, object.TenantId)
	if err != nil {
		log.Errorln("check rights failed, err:", err)
		return nil
	}

	out.Object = object.Object
	object.StorageClass, _ = GetNameFromTier(object.Tier, utils.OSTYPE_OPENSDS)
	return nil
}

func (s *s3Service) GetObject(ctx context.Context, req *pb.GetObjectInput, stream pb.S3_GetObjectStream) error {
	log.Infoln("GetObject is called in s3 service.")
	bucketName := req.Bucket
	objectName := req.Key
	offset := req.Offset
	length := req.Length

	var err error
	getObjRes := &pb.GetObjectResponse{}
	defer func() {
		getObjRes.ErrorCode = GetErrCode(err)
		stream.SendMsg(getObjRes)
	}()

	object, err := s.MetaStorage.GetObject(ctx, bucketName, objectName, "", true)
	if err != nil {
		log.Errorln("failed to get object info from meta storage. err:", err)
		return err
	}
	_, _, err = CheckRights(ctx, object.TenantId)
	if err != nil {
		return err
	}

	bucket, err := s.MetaStorage.GetBucket(ctx, bucketName, true)
	if err != nil {
		log.Errorln("failed to get bucket from meta storage. err:", err)
		return err
	}

	isAdmin, tenantId, err := util.GetCredentialFromCtx(ctx)
	if err != nil && isAdmin == false {
		log.Error("get tenant id failed")
		err = ErrInternalError
		return nil
	}

	if isAdmin == false {
		switch object.Acl.CannedAcl {
		case "bucket-owner-full-control":
			if bucket.TenantId != tenantId {
				err = ErrAccessDenied
				return err
			}
		default:
			if object.TenantId != tenantId {
				err = ErrAccessDenied
				return err
			}
		}
	}

	backendName := bucket.DefaultLocation
	if object.Location != "" {
		backendName = object.Location
	}
	// if this object has only one part
	backend, err := utils.GetBackend(ctx, s.backendClient, backendName)
	if err != nil {
		log.Errorln("unable to get backend. err:", err)
		return err
	}
	sd, err := driver.CreateStorageDriver(backend.Type, backend)
	if err != nil {
		log.Errorln("failed to create storage driver. err:", err)
		return err
	}
	log.Infof("get object offset %v, length %v", offset, length)
	reader, err := sd.Get(ctx, object.Object, offset, offset+length-1)
	if err != nil {
		log.Errorln("failed to get data. err:", err)
		return err
	}

	eof := false
	left := object.Size
	buf := make([]byte, ChunkSize)
	for !eof && left > 0 {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			log.Errorln("failed to read, err:", err)
			break
		}
		// From https://golang.org/pkg/io/, a Reader returning a non-zero number of bytes at the end of the input stream
		// may return either err == EOF or err == nil. The next Read should return 0, EOF.
		// If err is equal to io.EOF, a non-zero number of bytes may be returned.
		if err == io.EOF {
			log.Debugln("finished read")
			eof = true
		}
		// From https://golang.org/pkg/io/, there is the following statement.
		// Implementations of Read are discouraged from returning a zero byte count with a nil error, except when len(p) ==
		// 0. Callers should treat a return of 0 and nil as indicating that nothing happened; in particular it does not indicate EOF.
		// If n is equal 0, it indicate that there is no more data to read
		if n == 0 {
			log.Infoln("reader return zero bytes.")
			break
		}

		err = stream.Send(&pb.GetObjectResponse{ErrorCode: int32(ErrNoErr), Data: buf[0:n]})
		if err != nil {
			log.Infof("stream send error: %v\n", err)
			break
		}
		left -= int64(n)
	}

	log.Infoln("get object successfully")
	return err
}

func (s *s3Service) UpdateObjectMeta(ctx context.Context, in *pb.Object, out *pb.PutObjectResponse) error {
	log.Infoln("UpdateObjectMeta is called in s3 service.")
	var err error
	defer func() {
		out.ErrorCode = GetErrCode(err)
	}()

	object, err := s.MetaStorage.GetObject(ctx, in.BucketName, in.ObjectKey, "", true)
	if err != nil {
		log.Errorln("failed to get object info from meta storage. err:", err)
		return err
	}
	_, _, err = CheckRights(ctx, object.TenantId)
	if err != nil {
		return err
	}

	err = s.MetaStorage.UpdateObjectMeta(&meta.Object{Object: in})
	if err != nil {
		log.Errorf("failed to update object meta storage, err:", err)
		err = ErrInternalError
		return err
	}
	out.LastModified = in.LastModified
	out.Md5 = in.Etag
	out.VersionId = in.GetVersionId()

	return nil
}

func (s *s3Service) CopyObject(ctx context.Context, in *pb.CopyObjectRequest, out *pb.CopyObjectResponse) error {
	log.Infoln("CopyObject is called in s3 service.")
	var err error
	defer func() {
		out.ErrorCode = GetErrCode(err)
	}()

	srcBucketName := in.SrcBucketName
	srcObjectName := in.SrcObjectName
	targetBucketName := in.TargetBucketName
	targetObjectName := in.TargetObjectName
	srcBucket, err := s.MetaStorage.GetBucket(ctx, srcBucketName, true)
	if err != nil {
		log.Errorln("get bucket failed with err:", err)
		return err
	}
	srcObject, err := s.MetaStorage.GetObject(ctx, srcBucketName, srcObjectName, "", true)
	if err != nil {
		log.Errorln("failed to get object info from meta storage. err:", err)
		return err
	}
	_, tenantid, err := CheckRights(ctx, srcObject.TenantId)
	if err != nil {
		log.Errorf("no rights to access the source object[%s]\n", srcObject.ObjectKey)
		return nil
	}

	backendName := srcBucket.DefaultLocation
	if srcObject.Location != "" {
		backendName = srcObject.Location
	}
	srcBackend, err := utils.GetBackend(ctx, s.backendClient, backendName)
	if err != nil {
		log.Errorln("failed to get backend client with err:", err)
		return err
	}
	srcSd, err := driver.CreateStorageDriver(srcBackend.Type, srcBackend)
	if err != nil {
		log.Errorln("failed to create storage. err:", err)
		return err
	}

	targetBucket, err := s.MetaStorage.GetBucket(ctx, targetBucketName, true)
	if err != nil {
		log.Errorln("get bucket failed with err:", err)
		return err
	}
	targetBackendName := targetBucket.DefaultLocation
	targetBackend, err := utils.GetBackend(ctx, s.backendClient, targetBackendName)
	if err != nil {
		log.Errorln("failed to get backend client with err:", err)
		return err
	}
	targetSd, err := driver.CreateStorageDriver(targetBackend.Type, targetBackend)
	if err != nil {
		log.Errorln("failed to create storage. err:", err)
		return err
	}

	reader, err := srcSd.Get(ctx, srcObject.Object, 0, srcObject.Size)
	if err != nil {
		log.Errorln("failed to put data. err:", err)
		return err
	}
	limitedDataReader := io.LimitReader(reader, srcObject.Size)

	targetObject := &pb.Object{
		ObjectKey:  targetObjectName,
		BucketName: targetBucketName,
		Size:       srcObject.Size,
	}
	ctx = context.WithValue(ctx, dscommon.CONTEXT_KEY_SIZE, srcObject.Size)
	ctx = context.WithValue(ctx, dscommon.CONTEXT_KEY_MD5, srcObject.Etag)
	res, err := targetSd.Put(ctx, limitedDataReader, targetObject)
	if err != nil {
		log.Errorln("failed to put data. err:", err)
		return err
	}
	if res.Written < srcObject.Size {
		// TODO: delete incomplete object at backend
		log.Warnf("write objects, already written(%d), total size(%d)\n", res.Written, srcObject.Size)
		err = ErrIncompleteBody
		return err
	}

	targetObject.Etag = res.Etag
	targetObject.ObjectId = res.ObjectId
	targetObject.LastModified = time.Now().UTC().Unix()
	targetObject.ContentType = srcObject.ContentType
	targetObject.DeleteMarker = false
	targetObject.CustomAttributes = srcObject.CustomAttributes
	targetObject.Type = meta.ObjectTypeNormal
	targetObject.StorageMeta = res.Meta
	targetObject.Location = targetBackendName
	targetObject.TenantId = tenantid
	// we only support copy data with sse but not support copy data without sse right now
	targetObject.ServerSideEncryption = srcObject.ServerSideEncryption
	// TODO: delete old object

	err = s.MetaStorage.PutObject(ctx, &meta.Object{Object: targetObject}, nil, nil, true)
	if err != nil {
		log.Errorln("failed to put object meta. err:", err)
		// TODO: delete uncompleted object at backend
		err = ErrDBError
		return err
	}

	out.Md5 = res.Etag
	out.LastModified = targetObject.LastModified

	log.Infoln("Successfully copy object ", res.Written, " bytes.")
	return nil
}

func initTargeObject(ctx context.Context, in *pb.MoveObjectRequest, srcObject *pb.Object) (*pb.Object, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Error("get metadata from ctx failed.")
		return nil, ErrInternalError
	}

	targetObject := &pb.Object{
		ObjectKey:            srcObject.ObjectKey,
		BucketName:           srcObject.BucketName,
		ObjectId:             srcObject.ObjectId,
		Size:                 srcObject.Size,
		Etag:                 srcObject.Etag,
		Location:             srcObject.Location,
		Tier:                 srcObject.Tier,
		TenantId:             srcObject.TenantId,
		UserId:               srcObject.UserId,
		StorageMeta:          srcObject.StorageMeta,
		LastModified:         srcObject.LastModified,
		ContentType:          srcObject.ContentType,
		ServerSideEncryption: srcObject.ServerSideEncryption,
		Acl:                  srcObject.Acl,
		Type:                 srcObject.Type,
		DeleteMarker:         false,
		CustomAttributes:     md, /* TODO: only reserve http header attr*/
	}

	if in.TargetTier > 0 {
		targetObject.Tier = in.TargetTier
	}

	return targetObject, nil
}

// This is for lifecycle management.
func (s *s3Service) MoveObject(ctx context.Context, in *pb.MoveObjectRequest, out *pb.MoveObjectResponse) error {
	log.Infoln("MoveObject is called in s3 service.")

	err := s.checkMoveRequest(ctx, in)
	if err != nil {
		return err
	}

	srcObject, err := s.MetaStorage.GetObject(ctx, in.SrcBucket, in.SrcObject, in.SrcObjectVersion, true)
	if err != nil {
		log.Errorf("failed to get object[%s] of bucket[%s]. err:%v\n", in.SrcObject, in.SrcBucket, err)
		return err
	}

	targetObject, err := initTargeObject(ctx, in, srcObject.Object)
	if err != nil {
		log.Errorf("failed to get init target obejct. err:%v\n", err)
		return err
	}

	var srcSd, targetSd driver.StorageDriver
	var srcBucket, targetBucket *types.Bucket
	srcBucket, err = s.MetaStorage.GetBucket(ctx, in.SrcBucket, true)
	if err != nil {
		log.Errorf("get source bucket[%s] failed with err:%v", in.SrcBucket, err)
		return err
	}

	srcBackend, err := utils.GetBackend(ctx, s.backendClient, srcObject.Location)
	if err != nil {
		log.Errorln("failed to get backend client with err:", err)
		return err
	}
	srcSd, err = driver.CreateStorageDriver(srcBackend.Type, srcBackend)
	if err != nil {
		log.Errorln("failed to create storage. err:", err)
		return err
	}

	if in.MoveType == utils.MoveType_ChangeStorageTier {
		log.Infof("chagne storage class of %s\n", targetObject.ObjectKey)
		// just change storage tier
		targetBucket = srcBucket
		className, err := GetNameFromTier(in.TargetTier, srcBackend.Type)
		if err != nil {
			return ErrInternalError
		}
		err = srcSd.ChangeStorageClass(ctx, targetObject, &className)
		if err != nil {
			log.Errorf("change storage class of object[%s] failed, err:%v\n", targetObject.ObjectKey, err)
			return err
		}
		newObj := &meta.Object{Object: targetObject}
		err = s.MetaStorage.UpdateObject4Lifecycle(ctx, srcObject, newObj)
	} else {
		// need move data, get target location first
		if in.MoveType == utils.MoveType_ChangeLocation {
			targetBucket = srcBucket
			targetObject.Location = in.TargetLocation
			log.Infof("move %s cross backends, srcBackend=%s, targetBackend=%s, targetTier=%d\n",
				srcObject.ObjectKey, srcObject.Location, targetObject.Location, targetObject.Tier)
		} else { // MoveType_MoveCrossBuckets
			log.Infof("move %s from bucket[%s] to bucket[%s]\n", targetObject.ObjectKey, srcObject.BucketName,
				targetObject.BucketName)
			targetBucket, err = s.MetaStorage.GetBucket(ctx, in.TargetBucket, true)
			if err != nil {
				log.Errorf("get bucket[%s] failed with err:%v\n", in.TargetBucket, err)
				return err
			}
			targetObject.ObjectKey = in.TargetObject
			targetObject.Location = targetBucket.DefaultLocation
			targetObject.BucketName = targetBucket.Name
			log.Infof("move %s cross buckets, targetBucket=%s, targetBackend=%s, targetTier=%d\n",
				srcObject.ObjectKey, targetObject.BucketName, targetObject.Location, targetObject.Tier)
		}

		// get storage driver
		targetBackend, err := utils.GetBackend(ctx, s.backendClient, targetObject.Location)
		if err != nil {
			log.Errorln("failed to get backend client with err:", err)
			return err
		}
		targetSd, err = driver.CreateStorageDriver(targetBackend.Type, targetBackend)
		if err != nil {
			log.Errorln("failed to create storage. err:", err)
			return err
		}

		// copy data from one backend to another
		err = s.copyData(ctx, srcSd, targetSd, srcObject.Object, targetObject)
		if err != nil {
			log.Errorf("failed to copy object[%s], err:%v", srcObject.ObjectKey, err)
			return err
		}
		newObj := &meta.Object{Object: targetObject}
		if srcObject.Etag != targetObject.Etag {
			log.Errorf("data integrity check failed, etag of source object is %s, etag of target object is:%s\n",
				srcObject.Etag, targetObject.Etag)
			// if failed, delete target object
			s.cleanObject(ctx, newObj, targetSd)
			return err
		}
		out.Md5 = targetObject.Etag
		out.LastModified = targetObject.LastModified
		// update object meta data
		err = s.MetaStorage.UpdateObject4Lifecycle(ctx, srcObject, newObj)
		if err != nil {
			log.Errorln("failed to update meta data after copy, err:", err)
			// if failed, delete target object
			s.cleanObject(ctx, newObj, targetSd)
			return err
		}
		// delete source object
		s.cleanObject(ctx, srcObject, srcSd)
		log.Infof("delete source object[key=%s]\n", srcObject.ObjectKey)
	}

	log.Infoln("MoveObject is finished.")

	return nil
}

func (s *s3Service) copyData(ctx context.Context, srcSd, targetSd driver.StorageDriver, srcObj, targetObj *pb.Object) error {
	log.Infof("copy object data")
	reader, err := srcSd.Get(ctx, srcObj, 0, srcObj.Size)
	if err != nil {
		log.Errorln("failed to get data. err:", err)
		return err
	}
	limitedDataReader := io.LimitReader(reader, srcObj.Size)
	res, err := targetSd.Put(ctx, limitedDataReader, targetObj)
	if err != nil {
		log.Errorln("failed to put data. err:", err)
		return err
	}
	log.Infoln("Successfully copy ", res.Written, " bytes.")

	targetObj.Etag = res.Etag
	targetObj.ObjectId = res.ObjectId
	targetObj.StorageMeta = res.Meta

	return nil
}

func (s *s3Service) checkMoveRequest(ctx context.Context, in *pb.MoveObjectRequest) (err error) {
	log.Infoln("check copy request")

	if in.SrcBucket == "" || in.SrcObject == "" {
		log.Errorf("invalid copy source.")
		err = ErrInvalidCopySource
		return
	}

	switch in.MoveType {
	case utils.MoveType_ChangeStorageTier:
		if !validTier(in.TargetTier) {
			log.Error("cannot copy object to it's self.")
			err = ErrInvalidCopyDest
		}
	case utils.MoveType_ChangeLocation:
		if in.TargetLocation == "" {
			log.Errorf("no target lcoation provided for change location copy")
			err = ErrInvalidCopyDest
		}
		// in.TargetTier > 0 means need to change storage class
		if in.TargetTier > 0 {
			if !validTier(in.TargetTier) {
				log.Error("cannot copy object to it's self.")
				err = ErrInvalidCopyDest
			}
		}
	case utils.MoveType_MoveCrossBuckets:
	default:
		// copy cross buckets as default
		if in.TargetObject == "" || in.TargetBucket == "" {
			log.Errorf("invalid copy target")
			err = ErrInvalidCopyDest
			return
		}
		// in.TargetTier > 0 means need to change storage class
		if in.TargetTier > 0 {
			if !validTier(in.TargetTier) {
				log.Error("cannot copy object to it's self.")
				err = ErrInvalidCopyDest
			}
		}
	}

	log.Infof("MoveType:%d, srcObject:%v\n", in.MoveType, in.SrcObject)
	return
}

// When bucket versioning is Disabled/Enabled/Suspended, and request versionId is set/unset:
//
// |           |        with versionId        |                   without versionId                    |
// |-----------|------------------------------|--------------------------------------------------------|
// | Disabled  | error                        | remove object                                          |
// | Enabled   | remove corresponding version | add a delete marker                                    |
// | Suspended | remove corresponding version | remove null version object(if exists) and add a        |
// |           |                              | null version delete marker                             |
//
// See http://docs.aws.amazon.com/AmazonS3/latest/dev/Versioning.html
func (s *s3Service) DeleteObject(ctx context.Context, in *pb.DeleteObjectInput, out *pb.DeleteObjectOutput) error {
	log.Infoln("DeleteObject is called in s3 service, bucket:%s, key:%s, version:%s", in.Bucket, in.Key, in.VersioId)

	var err error
	defer func() {
		out.ErrorCode = GetErrCode(err)
	}()

	bucket, err := s.MetaStorage.GetBucket(ctx, in.Bucket, true)
	if err != nil {
		log.Errorln("get bucket failed with err:", err)
		return nil
	}

	object, err := s.MetaStorage.GetObject(ctx, in.Bucket, in.Key, in.VersioId, true)
	if err != nil {
		log.Errorln("failed to get object info from meta storage. err:", err)
		return err
	}
	isAdmin, tenantId, err := CheckRights(ctx, object.TenantId)
	if err != nil {
		log.Errorf("no rights to access the object[%s]\n", object.ObjectKey)
		return nil
	}

	// administrator can delete any resource
	if isAdmin == false {
		switch bucket.Acl.CannedAcl {
		case "public-read-write":
			break
		default:
			if bucket.TenantId != tenantId && tenantId != "" {
				log.Errorf("delete object failed: tenant[id=%s] has no access right.", tenantId)
				err = ErrBucketAccessForbidden
				return nil
			}
		} // TODO policy and fancy ACL
	}

	switch bucket.Versioning.Status {
	case utils.VersioningDisabled:
		err = s.removeObject(ctx, bucket, in.Key)
	case utils.VersioningEnabled:
		// TODO: versioning
		err = ErrInternalError
	case utils.VersioningSuspended:
		// TODO: versioning
		err = ErrInternalError
	default:
		log.Errorf("versioing of bucket[%s] is invalid:%s\n", bucket.Name, bucket.Versioning)
		err = ErrInternalError
	}

	// TODO: need to refresh cache if it is enabled

	return nil
}

func (s *s3Service) removeObject(ctx context.Context, bucket *meta.Bucket, objectKey string) error {
	log.Debugf("remove object[%s] from bucket[%s]\n", objectKey, bucket.Name)
	obj, err := s.MetaStorage.GetObject(ctx, bucket.Name, objectKey, "", true)
	if err == ErrNoSuchKey {
		return nil
	}
	if err != nil {
		log.Errorf("get object failed, err:%v\n", err)
		return ErrInternalError
	}

	backendName := bucket.DefaultLocation
	if obj.Location != "" {
		backendName = obj.Location
	}
	backend, err := utils.GetBackend(ctx, s.backendClient, backendName)
	if err != nil {
		log.Errorln("failed to get backend with err:", err)
		return err
	}
	sd, err := driver.CreateStorageDriver(backend.Type, backend)
	if err != nil {
		log.Errorln("failed to create storage, err:", err)
		return err
	}

	// mark object as deleted
	err = s.MetaStorage.MarkObjectAsDeleted(ctx, obj)
	if err != nil {
		log.Errorln("failed to mark object as deleted, err:", err)
		return err
	}

	// delete object data in backend
	err = sd.Delete(ctx, &pb.DeleteObjectInput{Bucket: bucket.Name, Key: objectKey, VersioId: obj.VersionId,
		ETag: obj.Etag, StorageMeta: obj.StorageMeta, ObjectId: obj.ObjectId})
	if err != nil {
		log.Errorf("failed to delete obejct[%s] from backend storage, err:%v\n", objectKey, err)
		return err
	} else {
		log.Infof("delete obejct[%s] from backend storage successfully.\n", objectKey)
	}

	// delete object meta data from database
	err = s.MetaStorage.DeleteObject(ctx, obj)
	if err != nil {
		log.Errorf("failed to delete obejct[%s] metadata, err:%v", objectKey, err)
	} else {
		log.Infof("delete obejct[%s] metadata successfully.", objectKey)
	}

	return err
}

func (s *s3Service) ListObjects(ctx context.Context, in *pb.ListObjectsRequest, out *pb.ListObjectsResponse) error {
	log.Infof("ListObject is called in s3 service, bucket is %s.\n", in.Bucket)
	var err error
	defer func() {
		out.ErrorCode = GetErrCode(err)
	}()
	// Check ACL
	bucket, err := s.MetaStorage.GetBucket(ctx, in.Bucket, true)
	if err != nil {
		log.Errorf("err:%v\n", err)
		return nil
	}

	isAdmin, tenantId, err := util.GetCredentialFromCtx(ctx)
	if err != nil {
		log.Error("get tenant id failed")
		err = ErrInternalError
		return nil
	}

	// administrator can get any resource
	if isAdmin == false {
		switch bucket.Acl.CannedAcl {
		case "public-read", "public-read-write":
			break
		case "authenticated-read":
			if tenantId == "" {
				log.Error("no tenant id provided")
				err = ErrBucketAccessForbidden
				return nil
			}
		default:
			if bucket.TenantId != tenantId {
				log.Errorf("tenantId(%s) does not much bucket.TenantId(%s)", tenantId, bucket.TenantId)
				err = ErrBucketAccessForbidden
				return nil
			}
		}
		// TODO validate user policy and ACL
	}

	retObjects, appendInfo, err := s.ListObjectsInternal(ctx, in)
	if appendInfo.Truncated && len(appendInfo.NextMarker) != 0 {
		out.NextMarker = appendInfo.NextMarker
	}
	if in.Version == constants.ListObjectsType2Int {
		out.NextMarker = util.Encrypt(out.NextMarker)
	}

	objects := make([]*pb.Object, 0, len(retObjects))
	for _, obj := range retObjects {
		object := pb.Object{
			LastModified: obj.LastModified,
			Etag:         obj.Etag,
			Size:         obj.Size,
			Tier:         obj.Tier,
			Location:     obj.Location,
			TenantId:     obj.TenantId,
			BucketName:   obj.BucketName,
			VersionId:    obj.VersionId,
		}
		if in.EncodingType != "" { // only support "url" encoding for now
			object.ObjectKey = url.QueryEscape(obj.ObjectKey)
		} else {
			object.ObjectKey = obj.ObjectKey
		}
		object.StorageClass, _ = GetNameFromTier(obj.Tier, utils.OSTYPE_OPENSDS)
		objects = append(objects, &object)
		log.Debugf("object:%+v\n", object)
	}
	out.Objects = objects
	out.Prefixes = appendInfo.Prefixes
	out.IsTruncated = appendInfo.Truncated

	if in.EncodingType != "" { // only support "url" encoding for now
		out.Prefixes = helper.Map(out.Prefixes, func(s string) string {
			return url.QueryEscape(s)
		})
		out.NextMarker = url.QueryEscape(out.NextMarker)
	}

	err = ErrNoErr
	return nil
}

func (s *s3Service) ListObjectsInternal(ctx context.Context, request *pb.ListObjectsRequest) (retObjects []*meta.Object,
	appendInfo utils.ListObjsAppendInfo, err error) {
	log.Infoln("Prefix:", request.Prefix, "Marker:", request.Marker, "MaxKeys:",
		request.MaxKeys, "Delimiter:", request.Delimiter, "Version:", request.Version,
		"keyMarker:", request.KeyMarker, "versionIdMarker:", request.VersionIdMarker)

	filt := make(map[string]string)
	if request.Versioned {
		filt[common.KMarker] = request.KeyMarker
		filt[common.KVerMarker] = request.VersionIdMarker
	} else if request.Version == constants.ListObjectsType2Int {
		if request.ContinuationToken != "" {
			var marker string
			marker, err = util.Decrypt(request.ContinuationToken)
			if err != nil {
				err = ErrInvalidContinuationToken
				return
			}
			filt[common.KMarker] = marker
		} else {
			filt[common.KMarker] = request.StartAfter
		}
	} else { // version 1
		filt[common.KMarker] = request.Marker
	}

	filt[common.KPrefix] = request.Prefix
	filt[common.KDelimiter] = request.Delimiter
	// currentlly, request.Filter only support filter by 'lastmodified' and 'tier'
	for k, v := range request.Filter {
		filt[k] = v
	}

	return s.MetaStorage.Db.ListObjects(ctx, request.Bucket, request.Versioned, int(request.MaxKeys), filt)
}

func (s *s3Service) cleanObject(ctx context.Context, object *Object, sd driver.StorageDriver) error {
	delInput := &pb.DeleteObjectInput{
		Bucket: object.BucketName, Key: object.ObjectKey, ObjectId: object.ObjectId,
		VersioId: object.VersionId, StorageMeta: object.StorageMeta,
	}

	err := sd.Delete(ctx, delInput)
	if err != nil {
		log.Warnf("clean object[%v] from backend failed, err:%v\n", err)
		ierr := s.MetaStorage.AddGcobjRecord(ctx, object)
		if ierr != nil {
			log.Warnf("add gc record failed, object:%v, err:%v\n", object, ierr)
		}
	}

	return err
}

func (s *s3Service) PutObjACL(ctx context.Context, in *pb.PutObjACLRequest, out *pb.BaseResponse) error {
	log.Info("PutObjACL is called in s3 service.")
	var err error
	defer func() {
		out.ErrorCode = GetErrCode(err)
	}()

	bucket, err := s.MetaStorage.GetBucket(ctx, in.ACLConfig.BucketName, true)
	if err != nil {
		log.Errorf("failed to get bucket meta. err: %v\n", err)
		return err
	}

	isAdmin, tenantId, err := util.GetCredentialFromCtx(ctx)
	if err != nil && isAdmin == false {
		log.Error("get tenant id failed")
		err = ErrInternalError
		return nil
	}

	// administrator can get any resource
	if isAdmin == false {
		switch bucket.Acl.CannedAcl {
		case "bucket-owner-full-control":
			if bucket.TenantId != tenantId {
				err = ErrAccessDenied
				return err
			}
		default:
			if bucket.TenantId != tenantId {
				err = ErrAccessDenied
				return err
			}
		}
		// TODO validate user policy and ACL
	}

	object, err := s.MetaStorage.GetObject(ctx, in.ACLConfig.BucketName, in.ACLConfig.ObjectKey, "", true)
	if err != nil {
		log.Errorln("failed to get object info from meta storage. err:", err)
		return err
	}
	object.Acl = &pb.Acl{CannedAcl: in.ACLConfig.CannedAcl}
	err = s.MetaStorage.UpdateObjectMeta(object)
	if err != nil {
		log.Infoln("failed to update object meta. err:", err)
		return err
	}

	log.Infoln("Put object acl successfully.")
	return nil
}
