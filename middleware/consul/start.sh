docker run -d --name consul \
	  -p 8500:8500 \
	    registry.1stcs.cn/zc/consul:1.15.4 agent -dev -ui -client=0.0.0.0