package jenkins

//Jenkins 安裝
//docker pull jenkins/jenkins
//mkdir jenkins_mount
//chown -R 1000:1000 /root/jenkins_mount
//docker run -d -p 10240:8080 -p 10241:50000 -v /root/jenkins_mount:/var/jenkins_home -v /etc/localtime:/etc/localtime --name jenkins jenkins/jenkins

//访问地址：
//http://172.16.71.130:10240/
//
//密码文件地址：
///Users/mac/libaobao/docker_mount/jenkins_mount/secrets
//密码：
//c3a80eebea5a4bbbb9f3d38e2dd527e8

//https://mirrors.tuna.tsinghua.edu.cn/jenkins/updates/update-center.json

//cd $WORKSPACE
//docker build . -t gin_swagger_jenkins
//docker tag localhost:5000/gin_swagger_jenkins
//docker push localhost:5000/gin_swagger_jenkins
