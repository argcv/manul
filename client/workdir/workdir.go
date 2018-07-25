package workdir

import (
	"github.com/argcv/go-argcvapis/app/manul/file"
	"github.com/argcv/manul/helpers"
	"github.com/argcv/webeh/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

type Workdir struct {
	Base   string
	Cwd    string
	isRoot bool
	parent *Workdir
}

func NewWorkdir(base string) *Workdir {
	base = path.Clean(base)
	dir, err := helpers.GetPathOfSelf()
	if err != nil {
		log.Errorf("Get Current GetCwd failed... %v", err)
		dir = "/tmp"
	}

	if base[0] != '/' {
		// relative path
		base = path.Join(dir, base)
	}

	os.MkdirAll(base, 0700)

	env := &Workdir{
		Base:   base,
		Cwd:    "/",
		isRoot: true,
	}
	env.parent = env
	return env
}

func (env *Workdir) Spawn() (sub *Workdir) {
	return &Workdir{
		Base:   env.Base,
		Cwd:    env.Cwd,
		isRoot: false,
		parent: env,
	}
}

func (env *Workdir) Rebase() (sub *Workdir) {
	return &Workdir{
		Base:   path.Join(env.Base, env.Cwd),
		Cwd:    "/",
		isRoot: false,
		parent: env,
	}
}

func (env *Workdir) GetRoot() *Workdir {
	if env.isRoot {
		return env
	} else {
		return env.parent.GetRoot()
	}
}

// Placeholder, close
func (env *Workdir) Close() {
	if env.isRoot {
	} else {
	}
}

func (env *Workdir) GetCwd() string {
	return path.Join(env.Base, env.Cwd)
}

func (env *Workdir) RemoveCwd() error {
	return os.RemoveAll(path.Join(env.Base, env.Cwd))
}

func (env *Workdir) Split() (dir, file string) {
	return path.Split(env.GetCwd())
}

func (env *Workdir) Path(filename string) string {
	return path.Join(env.Base, env.Cwd, filename)
}

func (env *Workdir) Remove(filename string) error {
	return os.RemoveAll(env.Path(filename))
}

func (env *Workdir) Goto(dir ...string) *Workdir {
	dirs := []string{
		env.Cwd,
	}
	dirs = append(dirs, dir...)
	env.Cwd = path.Join(dirs...)
	return env
}

func (env *Workdir) MkdirAll(filename string, perm os.FileMode) (err error) {
	dirTarget := env.Path(filename)
	if err = os.MkdirAll(dirTarget, perm); err != nil {
		log.Errorf("Create folder failed: %v", err)
		return
	} else {
		return nil
	}
}

func (env *Workdir) WriteFile(filename string, data []byte, perm os.FileMode) (err error) {
	tdir, tfile := path.Split(filename)
	dirTarget := path.Join(env.Base, env.Cwd, tdir)
	if err = os.MkdirAll(dirTarget, 0700); err != nil {
		log.Errorf("Create folder failed: %v", err)
		return
	} else {
		return ioutil.WriteFile(path.Join(dirTarget, tfile), data, perm)
	}
}

func (env *Workdir) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(env.Path(filename))
}

func (env *Workdir) Stat(filename string) (os.FileInfo, error) {
	return os.Stat(env.Path(filename))
}

func (env *Workdir) Exists(filename string) bool {
	if _, err := env.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func (env *Workdir) IsDir(filename string) bool {
	if mode, err := env.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return mode.IsDir()
	}
}

func (env *Workdir) IsFile(filename string) bool {
	if mode, err := env.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return !mode.IsDir()
	}
}

func (env *Workdir) ReadDir(filename string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(env.Path(filename))
}

func (env *Workdir) IterFiles(dir string, handler func(*file.File) error) (err error) {
	if !env.Exists(dir) {
		return errors.New("dir not exists")
	}
	if s, e := env.Stat(dir); e == nil {
		if s.IsDir() {
			if fl, err := env.ReadDir(dir); err != nil {
				log.Errorf("read dir %s(%s) failed: %v", dir, s.Name(), err)
				return err
			} else {
				for _, f := range fl {
					env.IterFiles(path.Join(dir, f.Name()), handler)
				}
			}
		} else {
			mode := s.Mode()
			size := s.Size()
			if data, e := env.ReadFile(dir); e != nil {
				log.Errorf("read file %s(%s) failed: %v", dir, s.Name(), err)
			} else {
				cpath, _ := path.Split(dir)
				meta := map[string]interface{}{
					"perm": mode,
				}
				fh := &file.File{
					Name: s.Name(),
					Path: cpath,
					Size: uint64(size),
					Meta: helpers.ToStruct(meta),
					Data: data,
				}
				handler(fh)
			}

		}

	}

	return nil
}
