#信链supervisor设计书(0.1.0版本)

##supervisor简介
信链supervisor角色是监控farmer节点的"记分员"、"监工",暂时由信链官方运营维护。
在区块链去中心化的理论框架下,适当的增加中心化的监察组件,有利于平台更快更健康的发展。

**注: _现阶段的supervisor为单节点设置, 如果配置成多节点可能会造成数据不一致_**

##supervisor系统架构
假装有图[我是图](http://www.baidu.com)

###farmer账号树
supervisor最重要的一个功能就是维持平台内所有farmer的选票(balance/token)信息、在线信息等。
选票是farmer为获得信链币,匿名通过teller为farmer代表进行投票的凭证(farmer所支持的farmer代表如果获得记账权,会将系统奖励的信链币的一部分分发给支持自己的farmer)。
一经消耗,则不能继续投票,也就不能获得信链币,只能通过向supervisor报告自己共享了多少区块,共享了多长时间,除此之外,并无其他方法获得选票。

**注: _teller进行匿名选举时,信息是从supervisor获取。对farmer来讲,选举过程并不参与,它只是贡献了一定区块空间和时间,就能获得一定量的信链币_**




