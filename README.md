#基于gin开启go自定义web基础组件构建篇

main.go中
* 1.配置初始化。2.日志初始化。3。models初始化。4.redis拓展再度封装并初始化。5.验证器组件初始化。
* 其实任何后期app中需要用到的服务，如果第三方包封装得不够简洁，我们便可在extend中做一层再度封装，那么业务逻辑中便可精简着使用
* 业务自身支撑的拓展包则在extend同级别目录建立