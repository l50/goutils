package utils

// func TestCloneRepo(t *testing.T) {
// 	targetDir := "/tmp"
// 	repo := "https://github.com/l50/helloworld.git"
// 	currentTime := time.Now()
// 	cloneLoc := filepath.Join(targetDir, fmt.Sprintf("helloworld-%s", currentTime.Format("2006-01-02-15-04-05")))

// 	cloned := CloneRepo(repo, cloneLoc)
// 	if !cloned {
// 		t.Fatalf("failed to clone %s", repo)
// 	}

// 	defer os.RemoveAll(cloneLoc)
// }
