## 引言

如果按照上个学期的内容，我们仅仅需要启动一个 go 进程即可启动完整的一个 web 服务，前面大家学习过微服务，也知道了现实的服务常常是集群的形式进行部署，单体服务我们如果做事故定位，瓶颈分析仅仅只需要简单的本地终端打日志、本地启个 pprof 端口这些手段来进行观察，然而在集群场景下，一次 API 接口调用往往涉及到很多个服务，此时想做错误定位、性能分析往往没办法依靠我们刚刚提到的手段了，那么怎么办？这就是我们这节课要讲的内容——**可观测性**。

虽然是可观测性的课，但是我觉得直接讲官方文档里面的 rolldice 还是太没劲了，而且我第一次看这玩意其实也听不懂，主要还是讲可观测性这些常见的组件以及他们的关系吧，实际上很多框架比如字节的 kitex、hertz，apache 的 dubbo-go 都有内置的一些包帮你们做好了内置的可观测性，一般来讲直接调包 + 部署一下可观测性的这些组件就好了。



## 正文

针对于单体架构的服务，我们常常有三种方式去做问题定位、性能分析：

1. 打日志
2. 链路，我们可能会在某个函数调用点打日志做关键的数据记录，比如耗时什么的。
3. pprof，其实这个我猜你们也用的不多，我之前主要是看堆内存分配，观察内存泄漏的。

   而在分布式集群的场景下，其实用的依旧是这一套，解决方法其实就是聚合分析，将集群产出的数据统一聚合到一个组件上进行观测，对应的也有三个名词，想必你们也听过：分布式日志、链路追踪、指标上报。

   不过值得一提的是链路追踪和上面引出来的**函数调用链**的概念并不一样，这里的链路追踪大多数情况都是以服务为单位的，当然其实也并不全是，也有部分服务内置的链路追踪组件往往会将其拆分成更细的粒度以方便观测。

   上面可能讲得有点抽象，来个具体的场景，以前我们是单体服务，比如一个电商，我们每次请求都是在这一个进程里面进行处理，最多就是调用一下数据库之类的，我们如果发现客户反馈支付之后订单显示没有支付成功，但是服务明显是正常运行的，那么我们就需要针对支付这个接口进行问题定位，往往我们就会通过日志以及看函数调用链的方式来进行分析。这当然很简单，因为只有一个进程。

   那么后来，我们的电商服务重构成了微服务架构，他并不是只在单一的进程里面跑了，就算是一个服务，它也有多个节点，此时我们原本的直接去看进程产出的日志就显得很原始了，因为你需要一个一个节点、一个一个服务去看对应的日志，效率非常低；同时在接口返回错误的时候，你不知道这次调用牵扯到了哪些服务，问题出在哪个服务上，此时的问题定位也显得很困难，于是才有了对应的微服务场景下的解决方案。



### 分布式日志

为啥需要分布式日志？你想想你如果有无数个容器，你看日志都必须要去 CLI 里面敲 docker logs 有多难受，其实对于这个问题只需要把所有服务产出的日志输入给一个日志存储中间件，然后通过一个仪表盘来看就行了，常见的解决方案比如 EFK（elasticsearch、filebeat、kibana）、ELK（elasticsearch、logstash、kibana），我个人只用过 ELK，其中 elasticsearch 是作为日志的存储中间件来运行的，kibana 则是仪表盘，logstash 以及 filebeat 都是来帮助收集进程产出的日志的一个服务，filebeat 我没用过，logstash 其实就是在服务器上读进程产出的日志文件，然后做过滤和一些处理，最终输出给 elasticsearch，总得来说并不算很复杂，如果你想试试，可以直接本地用 docker compose 部署下面的服务：

```yaml
 services:
 		elasticsearch:
        image: elasticsearch:8.17.0
        container_name: elasticsearch
        restart: always
        environment:
        - discovery.type=single-node
        - xpack.security.enabled=false
        - xpack.security.enrollment.enabled=false
        ports:
        - "9200:9200"
        - "9300:9300"
        volumes:
        - ./data/es_data:/usr/share/elasticsearch/data:rw
        networks:
            - app
    kibana:
        image: kibana:8.17.0
        container_name: kibana
        restart: always
        ports:
        - "5601:5601"
        environment:
        - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
        depends_on:
        - elasticsearch
        networks:
            - app
    logstash:
        image: logstash:8.17.0
        container_name: logstash
        restart: always
        volumes:
        - {你的日志产出目录}:/usr/share/logstash/logs/klog:ro
        - {你的日志产出目录}:/usr/share/logstash/logs/hlog:ro
        - {pipeline配置文件}:/usr/share/logstash/pipeline:rw
        depends_on:
        - elasticsearch
        networks:
            - app
            
```

这里给一个我用过的 pipeline 配置文件吧，需要根据实际情况改改：

```c
input {
  file {
    path => "/usr/share/logstash/logs/klog/*.log"
    start_position => "beginning"  # 从文件开头读取
    sincedb_path => "/dev/null"  # 不保存读取进度，重启容器时重新读取
  }

  file {
    path => "/usr/share/logstash/logs/hlog/*.log"
    start_position => "beginning"
    sincedb_path => "/dev/null"
  }
}

filter {
  json {
    source => "message"         # 来源字段
    # target => "parsed"        # 如果想把解析结果放到 parsed 子对象里，可取消注释
  }

#   （可选）删掉原始 message 字段，避免重复
  mutate {
    remove_field => [ "message" ]
  }
}

output {
  elasticsearch {
    hosts => ["http://elasticsearch:9200"]
    index => "logstash-%{+YYYY.MM.dd}"
  }

  stdout { codec => rubydebug }  # 打印日志到控制台
}
```

实际上也就部署几个中间件，除非你想当运维人员或者参与这几个中间件的研发，否则大多数情况我们都只是一个用户。



### 链路追踪

链路追踪以及后面要讲的指标的导出方式和分布式日志并不一样，大多数情况我们会使用 OpenTelemetry 的 SDK 编写数据的收集逻辑，但是 OpenTelemetry 本身并不提供数据的观测和聚合分析，它更像是定义了一个统一的标准，包括链路追踪、指标这些的数据格式统一进行适配的适配器，由 Otel SDK 导出的数据，首先是交给 Otel Collector，然后由 Otel Collector 统一进行分发给 prometheus 和 Jaeger 这种中间件。其实这种说法并不算准确，因为 Prometheus 一般是自己去 Collector 去拉数据的😃，所以这就是这群组件之间的关系，当然你其实用 Jaeger 或者 Prometheus 自己的 SDK 也是完全可以实现的，但是用 Otel SDK 可以适配大部分链路追踪和指标仪表盘的中间件，仅此而已。

这里再简单介绍一下链路追踪里面常见的一些概念，trace 和 span，一般来说，我们一次调用就会有一个 trace，表示这整个链路追踪的单位，但是可能不仅有一个 span，一个 span 往往是一个中间服务的调用，比如一次购买的链路是 网关->订单->库存->支付 这几个服务，那么整个链路就是一个 trace，而其中订单可能就是一个 span，当然一个 span 也可能有更加细的粒度。在分布式日志的聚合分析中，我们就可以通过一个 trace id 去分析一次调用里面产生的日志并且结合链路追踪里面提供的数据来进行具体分析，当然这也就对日志的输出内容有要求，比如 hertz 就提供了 ctxLogger ，它会从上下文中提取 traceid 并且打印出来，从而不需要我们自己再手动写输出 traceid 的逻辑。

它有一个什么作用呢？他可以将整条访问链路交给我们，然后帮助我们分析其每个中间服务的调用速度，同时它还可以携带一些关键的 attrbute，其实就是一些数据，帮助我们进行性能优化；除此之外，我们还可以通过错误传递定位到某次接口调用在整个链路中错误出在那个服务上，从而能够精准的定位错误。

链路追踪这块如果你去看官方文档的话，我感觉概念挺多的也挺杂的，虽然上面提到了 Otel SDK 的作用，但是我个人建议还是不要去深究怎么自己去用 Otel 的 SDK 去写插桩埋点逻辑，直接用框架帮你封装好的链路追踪包就行了。其实我们想要在服务里面做个链路追踪也很简单，部署个 Jaeger 和 Otel Collector，然后直接参考你使用的框架的文档直接接入即可。姑且给个我以前用过的 docker compose 吧：

```yaml
services:
		otel-collector:
        image: otel/opentelemetry-collector-contrib:latest
        container_name: otel-collector
        command: [ "--config=/etc/otel-collector-config.yaml" ]
        volumes:
        - otel-collector-config.yaml:/etc/otel-collector-config.yaml:rw
        ports:
        - "4317:4317"     # OTLP gRPC
        - "4318:4318"     # OTLP HTTP
        - "55679:55679"   # zPages (调试页面)
        - "8888:8888"     # Prometheus metrics
        - "8889:8889"     # Prometheus exporter
        - "13133:13133"   # 健康检查
        depends_on:
        - jaeger
        networks:
            - app
    jaeger:
        image: jaegertracing/all-in-one:latest
        container_name: jaeger
        environment:
        - COLLECTOR_OTLP_ENABLED=true
        ports:
        - "16686:16686"   # Jaeger UI
        - "14268:14268"   # Jaeger HTTP 接口
        - "14250:14250"   # gRPC 接口（Collector）
        - "6831:6831/udp" # Agent 接口
        networks:
            - app
```

当然还需要一个 Otel Collector 的配置文件：

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317

exporters:
  otlp:
    endpoint: jaeger:4317
    tls:
      insecure: true
  prometheus:
    endpoint: "0.0.0.0:8889"

processors:
  batch:

extensions:
  health_check:
  pprof:
    endpoint: :1888
  zpages:
    endpoint: :55679

service:
  extensions: [pprof, zpages, health_check]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlp]
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [prometheus]
```

其实这里还有包括了 prometheus 的配置，后面有用。



### 指标

指标这部分有 prometheus 就行了，再叠个 grafana 也行，都是仪表盘，通俗来讲就是通过在应用内部编写统计一些数据的代码，然后发送给 otel-collector 或者在本地暴露一个输出指标的服务器，然后让第三方仪表盘服务来抓取指标，最终在仪表盘上形成的一些图表，比如树状图、折线图什么的，总之也是上面日志、链路追踪那一套，只是数据格式以及最终的效果不一样而已。

不过 prometheus 除了展示指标还有一个关键的作用就是告警，比如某一段时间服务器的 cpu 使用率高达 90%、或者某个接口的响应时间 p99 特别大，那么这个时候就应该配置规则告警了，有一个组件叫 AlertManager，可以配置 im 机器人或者邮箱通知进行告警通知，还有一些比较关键的作用就是可以做性能分析，针对一些 cpu 使用率很高的应用或者某个服务的响应错误的次数很高，那么我们就可以针对性的进行调优或者 bugfix。

部署一下普罗米修斯：

```yml
services:
		prometheus:
        image: prom/prometheus:latest
        container_name: prometheus
        volumes:
        - ./data/prometheus:/etc/prometheus
        ports:
        - "9091:9090"
        networks:
            - app
    grafana:
        image: grafana/grafana:latest
        container_name: grafana
        ports:
        - "3000:3000"
        volumes:
        - ./data/grafana:/var/lib/grafana
        environment:
        - GF_SECURITY_ADMIN_USER=admin 
        - GF_SECURITY_ADMIN_PASSWORD=admin
        depends_on:
        - prometheus
        networks:
            - app
```

普罗米修斯默认是从其他服务器数据源上拉取指标来展示，给个配置 `prometheus.yml`，我们这里从上面的 otel-collector 上面拉取：

```yml
global:
  scrape_interval: 15s  # 每 15 秒抓一次

scrape_configs:
  - job_name: 'otel-collector'  # 自采自己的指标
    static_configs:
      - targets: ['otel-collector:8889']
```

grafana 是一个有很多数据源的仪表盘，在上面给了他的 docker-compose，数据源可以自己在 Web UI 上面进行配置，这里就不多赘述了，我实际上也没用过几次。



## 一点建议

这部分知识其实在面试的时候可能无关紧要，我们在大多数情况下虽然也不会直接去用原生的 Otel SDK 或者 Jaeger SDK 去写这些插桩埋点逻辑，而是去用框架提供的 adapter，几乎能够接近零代码就能享受到整套的可观测性的服务。

但是一个高可靠、可维护的系统是离不开可观测性的，我个人认为，你们至少需要会的是如何部署、如何在系统里面集成可观测性，并且让他可以运行起来，这是最基本的。然后可以进一步去了解这几个组件之间的关系，数据的流向，各自的工作方式等等，最后你感兴趣的话，可以去看框架内部是如何插桩埋点，来帮助你实现链路追踪、指标上报的，到这一步我觉得就完全足够了，所以其实课上没有讲太多内容，还得你们课后自己去研究，不过现在 AI 很厉害了，相信你们能学的很轻松。

