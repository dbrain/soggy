package soggy

import (
  "crypto/rand"
  "encoding/hex"
  "io"
  "fmt"
)

func UID() []byte {
  buf := make([]byte, 16)
  io.ReadFull(rand.Reader, buf)
  buf[6] = (buf[6] & 0x0f) | 0x40
  buf[8] = (buf[8] & 0x3f) | 0x80
  return buf
}

func UIDString() string {
  return hex.EncodeToString(UID())
}

func UUIDv4() string {
  b := UID();
  return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
