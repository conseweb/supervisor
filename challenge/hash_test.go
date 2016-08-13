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
	"gopkg.in/check.v1"
	pb "github.com/conseweb/supervisor/protos"
)

type HashTest struct {}

var _ = check.Suite(&HashTest{})

func (this *HashTest) TestMD5(c *check.C) {
	p := []byte("TestMD5")
	hash := HASH(pb.HashAlgo_MD5, p)
	c.Check(hash, check.Equals, "3877ff9691b889725d2fd53c43036726")
}

func (this *HashTest) TestSHA1(c *check.C) {
	p := []byte("TestSHA1")
	hash := HASH(pb.HashAlgo_SHA1, p)
	c.Check(hash, check.Equals, "a0e32680cecab6988c3492920b0dfa326b0e168a")
}

func (this *HashTest) TestSHA224(c *check.C) {
	p := []byte("TestSHA224")
	hash := HASH(pb.HashAlgo_SHA224, p)
	c.Check(hash, check.Equals, "897d051eeb53224c3ec8fc64547ef0d763e5e841713daa112f18298b")
}

func (this *HashTest) TestSHA256(c *check.C) {
	p := []byte("TestSHA256")
	hash := HASH(pb.HashAlgo_SHA256, p)
	c.Check(hash, check.Equals, "545c0b3c72adb5a1e4b85fd7394823de50b5e0b876edbf4caf7701bc1792f8a9")
}

func (this *HashTest) TestSHA384(c *check.C) {
	p := []byte("TestSHA384")
	hash := HASH(pb.HashAlgo_SHA384, p)
	c.Check(hash, check.Equals, "230251b20de7b16aa86258c237013047767f10c69355fbbf1c0a829cfee3c5811bf071ae6a653bd1224ee741c38fafee")
}

func (this *HashTest) TestSHA512(c *check.C) {
	p := []byte("TestSHA512")
	hash := HASH(pb.HashAlgo_SHA512, p)
	c.Check(hash, check.Equals, "461451a7f5ef9a466723b2a16cada5c3f1fc2145671c33f87c7c835ea2d296f6cddd020b293b4a3adf088b9f97268633daf77e499ec39c385852f72c94bcdd16")
}

func (this *HashTest) TestSHA3224(c *check.C) {
	p := []byte("TestSHA3224")
	hash := HASH(pb.HashAlgo_SHA3224, p)
	c.Check(hash, check.Equals, "0173cb75265655225be2c324810c2c0f95ba63da69c4a12288d6e040")
}

func (this *HashTest) TestSHA3256(c *check.C) {
	p := []byte("TestSHA3256")
	hash := HASH(pb.HashAlgo_SHA3256, p)
	c.Check(hash, check.Equals, "da91844ab1d75be323a5d6868408d2db6265d0746eec000eb190c1cc6a557aa7")
}

func (this *HashTest) TestSHA3384(c *check.C) {
	p := []byte("TestSHA3384")
	hash := HASH(pb.HashAlgo_SHA3384, p)
	c.Check(hash, check.Equals, "eb6f7f0d69992f0677f3bea604695fa3d4eb36bbe88611c578878bd808a1df24f4b300029e778f9a191e3fb0b82ba7b1")
}

func (this *HashTest) TestSHA3512(c *check.C) {
	p := []byte("TestSHA3512")
	hash := HASH(pb.HashAlgo_SHA3512, p)
	c.Check(hash, check.Equals, "2c9f782a9f9cfcdcbe98ed367a3fe62786acee59ed96fdda0fd30f9249ac91fe62d2102a284d21da18495d192767619cb8bd369ede8a29fa4d8d88091944eb5e")
}

func (this *HashTest) BenchmarkMD5(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_MD5, []byte("BenchmarkMD5"))
	}
}

func (this *HashTest) BenchmarkSHA1(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA1, []byte("BenchmarkSHA1"))
	}
}

func (this *HashTest) BenchmarkSHA224(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA224, []byte("BenchmarkSHA224"))
	}
}

func (this *HashTest) BenchmarkSHA256(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA256, []byte("BenchmarkSHA256"))
	}
}

func (this *HashTest) BenchmarkSHA384(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA384, []byte("BenchmarkSHA384"))
	}
}

func (this *HashTest) BenchmarkSHA512(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA512, []byte("BenchmarkSHA512"))
	}
}

func (this *HashTest) BenchmarkSHA3224(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA3224, []byte("BenchmarkSHA3224"))
	}
}

func (this *HashTest) BenchmarkSHA3256(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA3256, []byte("BenchmarkSHA3256"))
	}
}

func (this *HashTest) BenchmarkSHA3384(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA3384, []byte("BenchmarkSHA3384"))
	}
}

func (this *HashTest) BenchmarkSHA3512(c *check.C) {
	for i := 0; i < c.N; i++ {
		HASH(pb.HashAlgo_SHA3512, []byte("BenchmarkSHA3512"))
	}
}