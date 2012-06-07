package soggy

import (
  "testing"
)

func TestUID(t *testing.T) {
  uuid := UID();
  t.Logf("uuidBytes[%s]\n", uuid)
}

func TestUUIDv4(t *testing.T) {
  uuid := UUIDv4();
  t.Logf("uuidv4[%s]\n", uuid)
}

func TestUIDString(t *testing.T) {
  uuid := UIDString();
  t.Logf("uuidString[%s]\n", uuid)
}

func BenchmarkUIDString(b *testing.B) {
  seen := make(map[string]int, 1000)
  for i := 0; i < b.N; i++ {
     uid := UIDString();
     b.StopTimer()
     count := seen[uid]
     if count > 0 {
       b.Fatalf("duplicate uuid[%s] count %d", uid, count)
     }
     seen[uid] = count+1
     b.StartTimer()
  }
}
