package bls

import (
	"bytes"
	"strconv"
	"testing"
)

var unitN = 0

// Tests (for Benchmarks see below)
func testPre(t *testing.T) {
	t.Log("create secret key")
	m := []byte("this is a bls sample for go")
	var sec SecretKey
	sec.SetByCSPRNG()
	t.Log("sec:", sec.HexString())
	t.Log("create public key")
	pub := sec.GetPublicKey()
	t.Log("pub:", pub.HexString())
	sign := sec.Sign(m)
	t.Log("sign:", sign.HexString())
	if !sign.Verify(pub, m) {
		t.Error("Signature does not verify")
	}

	// How to make array of SecretKey
	{
		sec := make([]SecretKey, 3)
		for i := 0; i < len(sec); i++ {
			sec[i].SetByCSPRNG()
			t.Log("sec=", sec[i].HexString())
		}
	}
}

func testStringConversion(t *testing.T) {
	t.Log("testRecoverSecretKey")
	var sec SecretKey
	var s string
	if unitN == 6 {
		s = "16798108731015832284940804142231733909759579603404752749028378864165570215949"
	} else {
		s = "40804142231733909759579603404752749028378864165570215949"
	}
	err := sec.SetDecString(s)
	if err != nil {
		t.Fatal(err)
	}
	if s != sec.DecString() {
		t.Error("not equal")
	}
	s = sec.HexString()
	var sec2 SecretKey
	err = sec2.SetHexString(s)
	if err != nil {
		t.Fatal(err)
	}
	if !sec.IsEqual(&sec2) {
		t.Error("not equal")
	}
}

func testEachSign(t *testing.T, m []byte) ([]SecretKey, []PublicKey, []Sign) {
	n := 5
	secVec := make([]SecretKey, n)
	pubVec := make([]PublicKey, n)
	signVec := make([]Sign, n)

	for i := 0; i < n; i++ {
		err := secVec[i].SetLittleEndian([]byte{0, 0, 0, 0, 0, 0})
		if err != nil {
			t.Error(err)
		}

		pubVec[i] = *secVec[i].GetPublicKey()
		t.Logf("pubVec[%d]=%s\n", i, pubVec[i].HexString())

		if !pubVec[i].IsEqual(secVec[i].GetPublicKey()) {
			t.Errorf("Pubkey derivation does not match\n%s\n%s", pubVec[i].HexString(), secVec[i].GetPublicKey().HexString())
		}

		signVec[i] = *secVec[i].Sign(m)
		if !signVec[i].Verify(&pubVec[i], m) {
			t.Error("Pubkey derivation does not match")
		}
	}
	return secVec, pubVec, signVec
}

func testSign(t *testing.T) {
	m := []byte("testSign")
	t.Log(m)

	var sec0 SecretKey
	sec0.SetByCSPRNG()
	pub0 := sec0.GetPublicKey()
	s0 := sec0.Sign(m)
	if !s0.Verify(pub0, m) {
		t.Error("Signature does not verify")
	}
	testEachSign(t, m)
}

func testAdd(t *testing.T) {
	t.Log("testAdd")
	var sec1 SecretKey
	var sec2 SecretKey
	sec1.SetByCSPRNG()
	sec2.SetByCSPRNG()

	pub1 := sec1.GetPublicKey()
	pub2 := sec2.GetPublicKey()

	m := []byte("test test")
	sign1 := sec1.Sign(m)
	sign2 := sec2.Sign(m)

	t.Log("sign1    :", sign1.HexString())
	sign1.Add(sign2)
	t.Log("sign1 add:", sign1.HexString())
	pub1.Add(pub2)
	if !sign1.Verify(pub1, m) {
		t.Fail()
	}
}

func testSetValue(t *testing.T) {
	var sec1 SecretKey
	var sec2 SecretKey

	sec1.SetValue(1000)
	sec2.SetValue(1000)

	pub1 := sec1.GetPublicKey()
	pub2 := sec1.GetPublicKey()

	m := []byte("test test")
	sign1 := sec1.Sign(m)
	sign2 := sec2.Sign(m)

	if !sign1.Verify(pub2, m) {
		t.Errorf("the two secret keys derived different signatures," +
			" when they are supposed to be the same")
	}

	if !sign2.Verify(pub1, m) {
		t.Errorf("the two secret keys derived different signatures," +
			" when they are supposed to be the same")
	}

	if !bytes.Equal(sec1.LittleEndian(), sec2.LittleEndian()) {
		t.Errorf("two supposedly identical secret keys are not equal %v , %v", sec1.LittleEndian(), sec2.LittleEndian())
	}
}

func testData(t *testing.T) {
	t.Log("testData")
	var sec1, sec2 SecretKey
	sec1.SetByCSPRNG()
	b := sec1.LittleEndian()
	err := sec2.SetLittleEndian(b)
	if err != nil {
		t.Fatal(err)
	}
	if !sec1.IsEqual(&sec2) {
		t.Error("SecretKey not same")
	}
	pub1 := sec1.GetPublicKey()
	b = pub1.Serialize()
	var pub2 PublicKey
	err = pub2.Deserialize(b)
	if err != nil {
		t.Fatal(err)
	}
	if !pub1.IsEqual(&pub2) {
		t.Error("PublicKey not same")
	}
	m := []byte("doremi")
	sign1 := sec1.Sign(m)
	b = sign1.Serialize()
	var sign2 Sign
	err = sign2.Deserialize(b)
	if err != nil {
		t.Fatal(err)
	}
	if !sign1.IsEqual(&sign2) {
		t.Error("Sign not same")
	}
}

func testSerializeToHexStr(t *testing.T) {
	t.Log("testSerializeToHexStr")
	var sec1, sec2 SecretKey
	sec1.SetByCSPRNG()
	s := sec1.SerializeToHexStr()
	err := sec2.DeserializeHexStr(s)
	if err != nil {
		t.Fatal(err)
	}
	if !sec1.IsEqual(&sec2) {
		t.Error("SecretKey not same")
	}
	pub1 := sec1.GetPublicKey()
	s = pub1.SerializeToHexStr()
	var pub2 PublicKey
	err = pub2.DeserializeHexStr(s)
	if err != nil {
		t.Fatal(err)
	}
	if !pub1.IsEqual(&pub2) {
		t.Error("PublicKey not same")
	}
	m := []byte("doremi")
	sign1 := sec1.Sign(m)
	s = sign1.SerializeToHexStr()
	var sign2 Sign
	err = sign2.DeserializeHexStr(s)
	if err != nil {
		t.Fatal(err)
	}
	if !sign1.IsEqual(&sign2) {
		t.Error("Sign not same")
	}
}

func testOrder(t *testing.T, c int) {
	var curve string
	var field string
	if c == CurveFp254BNb {
		curve = "16798108731015832284940804142231733909759579603404752749028378864165570215949"
		field = "16798108731015832284940804142231733909889187121439069848933715426072753864723"
	} else if c == CurveFp382_1 {
		curve = "5540996953667913971058039301942914304734176495422447785042938606876043190415948413757785063597439175372845535461389"
		field = "5540996953667913971058039301942914304734176495422447785045292539108217242186829586959562222833658991069414454984723"
	} else if c == CurveFp382_2 {
		curve = "5541245505022739011583672869577435255026888277144126952448297309161979278754528049907713682488818304329661351460877"
		field = "5541245505022739011583672869577435255026888277144126952450651294188487038640194767986566260919128250811286032482323"
	} else if c == BLS12_381 {
		curve = "52435875175126190479447740508185965837690552500527637822603658699938581184513"
		field = "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787"
	} else {
		t.Fatal("bad c", c)
	}
	s := curveOrder()
	if s != curve {
		t.Errorf("bad curve order\n%s\n%s\n", s, curve)
	}
	s = fieldOrder()
	if s != field {
		t.Errorf("bad field order\n%s\n%s\n", s, field)
	}
}

func test(t *testing.T, c int) {
	err := initializeBLS(c)
	if err != nil {
		t.Fatal(err)
	}
	unitN = GetOpUnitSize()
	t.Logf("unitN=%d\n", unitN)
	testPre(t)
	testAdd(t)
	testSetValue(t)
	testSign(t)
	testData(t)
	testStringConversion(t)
	testOrder(t, c)
	testSerializeToHexStr(t)
}

func TestBLS(t *testing.T) {
	t.Logf("GetMaxOpUnitSize() = %d\n", GetMaxOpUnitSize())
	t.Log("CurveFp254BNb")
	test(t, CurveFp254BNb)
	if GetMaxOpUnitSize() == 6 {
		t.Log("CurveFp382_1")
		test(t, CurveFp382_1)
		t.Log("BLS12_381")
		test(t, BLS12_381)
	}
}

// Benchmarks

var curve = BLS12_381

//var curve = CurveFp254BNb

func BenchmarkPubkeyFromSeckey(b *testing.B) {
	b.StopTimer()
	err := initializeBLS(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.SetByCSPRNG()
		b.StartTimer()
		sec.GetPublicKey()
		b.StopTimer()
	}
}

func BenchmarkSigning(b *testing.B) {
	b.StopTimer()
	err := initializeBLS(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.SetByCSPRNG()
		b.StartTimer()
		sec.Sign([]byte(strconv.Itoa(n)))
		b.StopTimer()
	}
}

func BenchmarkValidation(b *testing.B) {
	b.StopTimer()
	err := initializeBLS(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.SetByCSPRNG()
		pub := sec.GetPublicKey()
		m := []byte(strconv.Itoa(n))
		sig := sec.Sign(m)
		b.StartTimer()
		sig.Verify(pub, m)
		b.StopTimer()
	}
}
