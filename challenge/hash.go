/*
Copyright Mojing Inc. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package challenge

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	pb "github.com/conseweb/common/protos"
	"golang.org/x/crypto/sha3"
	"hash"
)

func HASH(hashAlgo pb.HashAlgo, p []byte) string {
	var h hash.Hash

	switch hashAlgo {
	case pb.HashAlgo_MD5:
		h = md5.New()
	case pb.HashAlgo_SHA1:
		h = sha1.New()
	case pb.HashAlgo_SHA224:
		h = sha256.New224()
	case pb.HashAlgo_SHA256:
		h = sha256.New()
	case pb.HashAlgo_SHA384:
		h = sha512.New384()
	case pb.HashAlgo_SHA512:
		h = sha512.New()
	case pb.HashAlgo_SHA3224:
		h = sha3.New224()
	case pb.HashAlgo_SHA3256:
		h = sha3.New256()
	case pb.HashAlgo_SHA3384:
		h = sha3.New384()
	case pb.HashAlgo_SHA3512:
		h = sha3.New512()
	default:
		h = sha256.New()
	}

	h.Write(p)

	return fmt.Sprintf("%x", h.Sum(nil))
}
