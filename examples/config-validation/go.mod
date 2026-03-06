module github.com/keksclan/goStartyUpy/examples/config-validation

go 1.24

require (
	github.com/keksclan/goConfy v0.1.0
	github.com/keksclan/goStartyUpy v0.2.0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect

replace github.com/keksclan/goStartyUpy => ../..
