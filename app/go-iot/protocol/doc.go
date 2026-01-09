package protocol

// codec 用于对设备上报和云端指令下发进行编解码处理，用户可以实现Codec接口来自定义编解码器， 也可以使用插件的方式来支持扩展，

//对设备上报时， 基于产品id,  获取codec
