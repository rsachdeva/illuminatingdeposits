package userauthn

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// No mocking; just passwing parameters to trigger error
func TestPasswordNotMatchingNoMock(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("hello"), bcrypt.DefaultCost)
	pv := PasswordVerifier{hashedPassword: hash, password: "bye"}
	err := PasswordMatch(pv)
	require.NotNil(t, err, "should be an error for password mismatch")
	require.Regexp(t, regexp.MustCompile("combination not found"), err)
}

func TestPasswordMatchingNoMock(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("hello"), bcrypt.DefaultCost)
	l := len([]byte("hello"))
	fmt.Printf("l is %v\n", l)
	pv := PasswordVerifier{hashedPassword: hash, password: "hello"}
	err := PasswordMatch(pv)
	require.Nil(t, err)
}

// With mocking for demo purposes only to compare in case we had to mock
// not really needed here as just a library package
// https://github.com/stretchr/testify#mock-package
// Mock bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) as an example of mocking
// Mock the behavior defined in PasswordVeriferInterface interface
type MockedPasswordVerifier struct {
	mock.Mock
}

func (mpv *MockedPasswordVerifier) CompareHashAndPassword() error {
	args := mpv.Called()
	return args.Error(0)
}

func TestPasswordNotMatching(t *testing.T) {
	mpv := MockedPasswordVerifier{}

	// setup expectations
	mpv.On("CompareHashAndPassword").Return(errors.New("matching issue"))

	// call the code we are testing
	err := PasswordMatch(&mpv)
	require.NotNil(t, err, "should be an error for password mismatch")

}

func TestPasswordMatches(t *testing.T) {
	mpv := MockedPasswordVerifier{}

	// setup expectations
	mpv.On("CompareHashAndPassword").Return(nil)

	// call the code we are testing
	err := PasswordMatch(&mpv)

	require.Nil(t, err)
}
