//+build wireinject

package wiref

//安装wire，
//go get github.com/google/wire/cmd/wire
//上面的命令会在$GOPATH/bin中生成一个可执行程序wire，这就是代码生成器。我个人习惯把$GOPATH/bin加入系统环境变量$PATH中，所以可直接在命令行中执行wire命令。

//当前目录下，执行wire命令生产 wire_gen.go文件， 文件第一行 //+build wireinject，
//wire命令会针对包（wiref）生成一个统一的 wire_gen.go 文件， wire文件修改后，在包目录（wiref）下执行 go generate 会重新生成。

import (
	"github.com/google/wire"
	"go-lab/app/go_wire/internal/model"
)

func InitMission(name string) model.Mission {
	wire.Build(model.NewMonster, model.NewPlayer, model.NewMission)
	return model.Mission{}
}

//有时候，我们需要为某个类型绑定一个值，而不想依赖构造器每次都创建一个新的值。
//有些类型天生就是单例，例如配置，数据库对象（sql.DB）。这时我们可以使用wire.Value绑定值，使用wire.InterfaceValue绑定接口。
//例如下面，我们的怪兽一直是一个Kitty，我们就不用每次都去创建它了，直接绑定这个值就 ok 了：
var kitty = model.Monster{Name: "kitty"}

func InitMissionWithValue(name string) model.Mission {
	wire.Build(wire.Value(kitty), model.NewPlayer, model.NewMission)
	return model.Mission{}
}

//注意：InitMissionWithParam， 存在多个相同类型的参数时，
//最后使用别名类型进行区分， 如使用PlayerParam， MonsterParam 作为string的别名类型，已方便wire进行区分，否则wire生成文件的时候会报错
func InitMissionWithParam(p model.PlayerParam, m model.MonsterParam) model.Mission {
	wire.Build(model.NewMonsterWithParam, model.NewPlayerWithParam, model.NewMission)
	return model.Mission{}
}

var monsterPlayerSet = wire.NewSet(model.NewMonsterWithParam, model.NewPlayerWithParam)

func InitMissionWithParamSet(p model.PlayerParam, m model.MonsterParam) model.Mission {
	wire.Build(monsterPlayerSet, model.NewMission)
	return model.Mission{}
}

//总结
//wire是 Google 开源的一个依赖注入工具。它是一个代码生成器，并不是一个框架。
//我们只需要在一个特殊的go文件中告诉wire类型之间的依赖关系，它会自动帮我们生成代码，帮助我们创建指定类型的对象，并组装它的依赖。
//不像其他开源的依赖注入工具（uber和facebook的），不使用反射，所以不会存在性能问题。
