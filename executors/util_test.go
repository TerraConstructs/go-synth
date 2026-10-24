package executors

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"github.com/terraconstructs/go-synth/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_copyDir(t *testing.T) {
	logger := getPrettyLogger()
	// create test filesystem
	files := []string{
		"foo/test.txt",
		"foo/subdir1/test.txt",
		"foo/subdir2/test.txt",
	}
	fs := afero.NewMemMapFs()
	for _, path := range files {
		err := afero.WriteFile(fs, path, []byte("Hello, World!"), 0644)
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		fromDir       string
		toDir         string
		options       models.CopyOptions
		expectedFiles map[string]bool
	}{
		{
			name:    "Default should copy everything",
			fromDir: "foo",
			toDir:   "bar",
			options: models.CopyOptions{},
			expectedFiles: map[string]bool{
				"bar/test.txt":         true,
				"bar/subdir1/test.txt": true,
				"bar/subdir2/test.txt": true,
			},
		},
		{
			name:    "ignore all files except test.txt",
			fromDir: "foo",
			toDir:   "bar",
			options: models.CopyOptions{
				IgnorePatterns: []string{"*/*"},
				AllowPatterns:  []string{"test.txt"},
			},
			expectedFiles: map[string]bool{
				"bar/test.txt":         true,
				"bar/subdir1/test.txt": false,
				"bar/subdir2/test.txt": false,
			},
		},
		{
			name:    "skip directory",
			fromDir: "foo",
			toDir:   "bar",
			options: models.CopyOptions{
				SkipDirs: []string{"subdir1"},
			},
			expectedFiles: map[string]bool{
				"bar/test.txt":         true,
				"bar/subdir1/test.txt": false,
				"bar/subdir2/test.txt": true,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			destFs := afero.NewMemMapFs()
			err := copyDir(logger, tc.fromDir, tc.toDir, fs, destFs, tc.options)
			require.NoError(t, err)
			failed := false
			for path, shouldExist := range tc.expectedFiles {
				exists := fileExists(destFs, path)
				if exists != shouldExist {
					t.Logf("expected file %s: got %t - expected %t", path, exists, shouldExist)
					failed = true
				} else {
					t.Logf("expected file %s: got %t", path, exists)
				}
			}
			if failed {
				t.Fail()
			}
		})
	}
}

func fileExists(fs afero.Fs, path string) bool {
	_, err := fs.Stat(path)
	return err == nil
}

func getPrettyLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Optional: colorize the log level
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // Optional: use ISO8601 time format
	config.Encoding = "console"
	logger, _ := config.Build()
	return logger
}

func Test_newTempFs(t *testing.T) {
	t.Run("honors GO_SYNTH_TMPDIR", func(t *testing.T) {
		base := filepath.Join(t.TempDir(), "nested")
		t.Setenv(TempDirEnvVar, base)

		fs, dir, err := newTempFs("go-synth-test")
		require.NoError(t, err)
		require.Equal(t, base, filepath.Dir(dir))

		require.NoError(t, afero.WriteFile(fs, "probe.txt", []byte("hi"), 0644))
		_, err = os.Stat(filepath.Join(dir, "probe.txt"))
		require.NoError(t, err)
	})

	t.Run("defaults to user cache dir on darwin", func(t *testing.T) {
		if runtime.GOOS != "darwin" {
			t.Skip("darwin only")
		}
		t.Setenv(TempDirEnvVar, "")
		cacheDir, err := os.UserCacheDir()
		require.NoError(t, err)

		_, dir, err := newTempFs("go-synth-test")
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(dir) })
		require.Equal(t, filepath.Join(cacheDir, "go-synth"), filepath.Dir(dir))
	})
}
