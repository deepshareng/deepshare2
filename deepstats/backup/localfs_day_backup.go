package backup

import (
	"bufio"
	"io"
	"os"
	"path/filepath"

	pb "github.com/MISingularity/deepshare2/deepstats/backup/proto"
	"github.com/golang/protobuf/proto"
)

// delimiter : "#"
type LocalfsBackup struct {
	path            string
	currentFilename string
	file            io.WriteCloser
}

var filechan chan string

func filename(event pb.Event) string {
	return convertTime(event.Timestamp) + "-event.dat"
}

func newLocalWriter(filename string) (io.WriteCloser, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
}

func (l *LocalfsBackup) Insert(event pb.Event) error {
	if filename(event) != l.currentFilename {
		if l.file != nil {
			l.file.Close()
		}
		l.currentFilename = filename(event)
		var err error
		l.file, err = newLocalWriter(filepath.Join(l.path, l.currentFilename))
		if err != nil {
			return err
		}
	}
	result, err := proto.Marshal(&event)
	if err != nil {
		return err
	}
	result = append(result, []byte("#")...)
	_, err = l.file.Write(result)
	if err != nil {
		return err
	}
	return nil
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		filechan <- path
	}
	return nil
}

func traverse(path string) {
	filepath.Walk(path, visit)
	close(filechan)
}

func (l *LocalfsBackup) RetriveAllEvents() ([]pb.Event, error) {
	filechan = make(chan string)
	go traverse(l.path)
	events := make([]pb.Event, 0)
	for {
		select {
		case path := <-filechan:
			if path == "" {
				return events, nil
			}
			r, err := os.OpenFile(path, os.O_RDONLY, 0)
			if err != nil {
				return []pb.Event{}, err
			}
			bufr := bufio.NewReader(r)
			for {
				data, err := bufr.ReadBytes('#')
				if err != nil {
					if err == io.EOF {
						break
					}
					return []pb.Event{}, err
				}
				data = data[:len(data)-1]
				event := pb.Event{}
				err = proto.Unmarshal(data, &event)
				if err != nil {
					return []pb.Event{}, err
				}
				events = append(events, event)
			}
		}
	}
}

func NewLocalFSBackupService(path string) (BackupService, error) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return nil, err
	}
	return &LocalfsBackup{path: path}, nil
}
