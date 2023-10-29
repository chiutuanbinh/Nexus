package storage

type CommitLogger interface {
	Write(key []byte, value []byte) error
	Load(readFunc func([]byte, []byte) error) error
}

type SimpleCommitLogger struct {
	
}

func CreateCommitLog(commitFileName string) CommitLogger {

	return &SimpleCommitLogger{CommitFileName: commitFileName}
}

func (c *SimpleCommitLogger) Write(key []byte, value []byte) error {

	return nil
}

func (c *SimpleCommitLogger) Load(readFunc func([]byte, []byte) error) error {
	return nil
}
