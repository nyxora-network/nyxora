package packager

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Packager struct {
	baseDir string
	cacheDir string
}

func NewPackager(baseDir string) *Packager {
	cacheDir := filepath.Join(baseDir, "cache")
	_ = os.MkdirAll(cacheDir, 0755)
	return &Packager{
		baseDir: baseDir,
		cacheDir: cacheDir,
	}
}

func (p *Packager) Pack(name, sourceDir string) error {
	log.Printf("[packager] packing %s from %s", name, sourceDir)

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("source directory %s does not exist", sourceDir)
	}

	destPath := filepath.Join(p.baseDir, "tunnels", name+".tar.gz")
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("create tunnels dir: %w", err)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create archive: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		}

		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("write header: %w", err)
		}

		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}

		return nil
	})
}

func (p *Packager) Unpack(name, destDir string) error {
	archivePath := filepath.Join(p.baseDir, "tunnels", name+".tar.gz")
	log.Printf("[packager] unpacking %s to %s", archivePath, destDir)

	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return fmt.Errorf("archive %s not found", archivePath)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("create dest dir: %w", err)
	}

	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar header: %w", err)
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("create dir %s: %w", target, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("create parent dir: %w", err)
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("create file %s: %w", target, err)
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("write file %s: %w", target, err)
			}
			f.Close()
		}
	}

	log.Printf("[packager] unpacked %s to %s", name, destDir)
	return nil
}

func (p *Packager) List() ([]string, error) {
	tunnelsDir := filepath.Join(p.baseDir, "tunnels")
	if _, err := os.Stat(tunnelsDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(tunnelsDir)
	if err != nil {
		return nil, fmt.Errorf("read tunnels dir: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tar.gz") {
			names = append(names, strings.TrimSuffix(entry.Name(), ".tar.gz"))
		}
	}
	return names, nil
}

func (p *Packager) Exists(name string) bool {
	archivePath := filepath.Join(p.baseDir, "tunnels", name+".tar.gz")
	_, err := os.Stat(archivePath)
	return err == nil
}

func (p *Packager) Remove(name string) error {
	archivePath := filepath.Join(p.baseDir, "tunnels", name+".tar.gz")
	return os.Remove(archivePath)
}
