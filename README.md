# snowFlake 雪花算法
snowflake ID 算法是 twitter 使用的唯一 ID 生成算法，为了满足 Twitter 每秒上万条消息的请求，使每条消息有唯一、有一定顺序的 ID ，且支持分布式生成。

![image](https://github.com/charizer/blog/blob/master/images/snowflake.jpg)

# 构成

snowflake ID 的结构是一个 64 bit 的 int 型数据。如下图所示：

##### 第1位bit：
二进制中最高位为1的都是负数，但是我们所需要的id应该都是整数，所以这里最高位应该为0
##### 后面的41位bit：
用来记录生成id时的毫秒时间戳，这里毫秒只用来表示正整数(计算机中正整数包含0)，所以可以表示的数值范围是0至2^41 - 1
##### 再后面的10位bit：
用来记录工作机器的id
10位bit可以表示的最大正整数为:0~ 2^10-1, 所以当前规则允许分布式最大节点数为1024个节点 我们可以根据业务需求来具体分配worker数和每台机器1毫秒可生成的id序号number数
##### 最后的12位：
用来表示单台机器每毫秒生成的id序号
12位bit可以表示的最大正整数为:0~2^12，即可用0、1、2、3...4095这4096(注意是从0开始计算)个数字来表示1毫秒内机器生成的序号(这个算法限定单台机器1毫秒内最多生成4096个id，超出则等待下一毫秒再生成)

### 基本数据结构：
```
const (
	workerBits uint8 = 10						// 节点数
	seqBits uint8 = 12						    // 1毫秒内可生成的id序号的二进制位数
	workerMax int64 = -1 ^ (-1 << workerBits)   // 节点ID的最大值，用于防止溢出
	seqMax int64 = -1 ^ (-1 << sewqBits)        // 同上，用来表示生成id序号的最大值
	timeShift uint8 = workerBits + seqBits      // 时间戳向左的偏移量
	workerShift uint8 = seqBits                 // 节点ID向左的偏移量
	epoch int64 = 1567906170596					// 开始运行时间
)

type Worker struct {
	// 添加互斥锁 确保并发安全
	mu        sync.Mutex
	// 记录时间戳
	timestamp int64
	// 该节点的ID
	workerId  int64
	// 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
	seq       int64
}

```

### 算法描述：
![image](https://github.com/charizer/blog/blob/master/images/snowflake_flow.jpg)

### 使用示例：
```
worker, err := NewWorker(1)
if err != nil {
	log.Errorf("new worker err:%v",err)
	return
}
log.Info("---worker 1---")
id := worker.Next()
log.Infof("id:%d",id)
id = worker.Next()
log.Infof("id:%d",id)
	
```

