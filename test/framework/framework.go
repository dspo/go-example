package framework

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"go.uber.org/zap"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"gitee.com/huajinet/go-example/pkg/database"
)

const (
	Namespace = "go-example-e2e"
)

var (
	//go:embed manifests/app0.yaml
	_app0Spec string

	//go:embed manifests/mysql.yaml
	_mysqlSpec string
)

var (
	_f *Framework
)

func init() {
	var err error
	if _f == nil {
		_f, err = newFramework()
		if err != nil {
			log.Fatalln("failed to new framework")
		}
	}
}

type Framework struct {
	scaffold *KubernetesScaffold
	db       *gorm.DB
}

func GetFramework() *Framework {
	if _f == nil {
		log.Fatalln("framework is not initialized")
	}
	return _f
}

func (f *Framework) DB() *gorm.DB {
	return f.db
}

func (f *Framework) DeployComponents(t *testing.T) {
	f.deployDatabase(t)
	f.initDatabase(t)
	f.connectDatabase(t)
	f.deployApp0(t)
	f.deployApp1(t)
}

func (f *Framework) deployApp0(t *testing.T) {
	t.Log("it is going to deploy app0")
	err := k8s.KubectlApplyFromStringE(t, f.scaffold.kubectlOptions, _app0Spec)
	assertNilErr(t, err)

	err = f.ensureServiceWithTimeout(t.Context(), "app0", f.scaffold.kubectlOptions.Namespace, 1, 30)
	assertNilErr(t, err)
}

func (f *Framework) deployApp1(t *testing.T) {
	// not implement yet
}

func (f *Framework) deployDatabase(t *testing.T) {
	t.Log("it is going to deploy MySQL")
	err := k8s.KubectlApplyFromStringE(t, f.scaffold.kubectlOptions, _mysqlSpec)
	assertNilErr(t, err)

	err = f.ensureServiceWithTimeout(t.Context(), "mysql", f.scaffold.kubectlOptions.Namespace, 1, 60)
	assertNilErr(t, err)
}

func (f *Framework) initDatabase(t *testing.T) {
	t.Log("it is going to init MySQL")
	db, err := sql.Open("mysql", "root:changeme@tcp(mysql:3306)/")
	assertNilErr(t, err)

	defer func() { _ = db.Close() }()

	err = db.Ping()
	assertNilErr(t, err)

	// try to drop database
	_, err = db.Exec("DROP DATABASE IF EXISTS `go_dev`")
	assertNilErr(t, err)

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `go_dev`")
	assertNilErr(t, err)
}

func (f *Framework) connectDatabase(t *testing.T) {
	f.db = database.NewGorm()
}

func (f *Framework) ensureServiceWithTimeout(ctx context.Context, name, namespace string, desiredEndpoints, timeout int) error {
	backoff := wait.Backoff{
		Duration: 6 * time.Second,
		Factor:   1,
		Steps:    timeout / 6,
	}
	var lastErr error
	condFunc := func() (bool, error) {
		ep, err := f.scaffold.clientset.CoreV1().Endpoints(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			lastErr = err
			log.Println("ERROR: failed to list endpoints",
				zap.String("service", name),
				zap.Error(err),
			)
			return false, nil
		}
		count := 0
		for _, ss := range ep.Subsets {
			count += len(ss.Addresses)
		}
		if count == desiredEndpoints {
			return true, nil
		}
		log.Println("INFO: endpoints count mismatch",
			zap.String("service", name),
			zap.Any("ep", ep),
			zap.Int("expected", desiredEndpoints),
			zap.Int("actual", count),
		)
		lastErr = fmt.Errorf("expected endpoints: %d but seen %d", desiredEndpoints, count)
		return false, nil
	}

	err := wait.ExponentialBackoff(backoff, condFunc)
	if err != nil {
		return lastErr
	}
	return nil
}

func assertNilErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("failed to for err: %v", err)
	}
}

func newFramework() (*Framework, error) {
	_f = new(Framework)

	var (
		err error
	)

	_f.scaffold, err = NewKubernetesScaffold(KubectlOptions{
		ContextName: "",
		ConfigPath:  "",
		Namespace:   Namespace,
	})
	if err != nil {
		return nil, err
	}

	return _f, nil
}
