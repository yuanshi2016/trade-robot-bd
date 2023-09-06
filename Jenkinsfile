pipeline{
    agent any
    environment{
        HARBOR_HOST='harbor.yuanshi01.com:30687'
        HARBOR_ADDR='harbor.yuanshi01.com:30687/trade'
        K8S_NAMESPACE='develop'
    }
    parameters {
//         string(name: 'PROJECT_NAME', defaultValue: 'common-svc', description: 'project name,same as the name ofdocker container')
        string(name: 'CONTAINER_VERSION', defaultValue: '', description: 'docker container version number, SET when major version number changed')
        // booleanParam(name: 'DEPLOYMENT_K8S', defaultValue: false, description: 'release deployment k8s')
    }
    stages {
        stage('Initial') {
            steps{
                script {
                        git branch: 'develop-kratos', credentialsId: '68dea4dd-3625-4d84-90aa-daf45f3a391a', url: 'git@github.com:yuanshi2016/trade-robot-bd.git'
                        PROJECT_NAME = env.JOB_NAME.split('/')[1]
                        DOCKER_IMAGE = env.JOB_NAME.split('/')[1]
                         APP_NAME = env.JOB_NAME.split('/')[1]
                         if (APP_NAME ==~ /^api-.*/) {
                             env.TARGET_PATH = "./${APP_NAME}"
                         } else {
                            env.TARGET_PATH = "./app/${APP_NAME}"
                         }
                          // 脚本式创建一个环境变量
                        if (params.CONTAINER_VERSION == '') {
                                env.APP_VERSION=PROJECT_NAME
//                             env.APP_VERSION = sh(returnStdout:true,script:"/home/jenkins-build-tools gen -p ${params.PROJECT_NAME}").trim()
                        }else {
                            env.APP_VERSION ="${params.CONTAINER_VERSION}-alpha"
                        }
                        sh "echo ${env.APP_VERSION}"
                    }
                }
        }
        stage("Docker Build") {
            when {
                allOf {
                    expression { env.APP_VERSION != null }
                }
            }
            steps("Start Build") {
                sh "docker login -u admin -p Harbor12345 ${HARBOR_HOST}"
                sh "docker build --build-arg TARGET_PATH=${TARGET_PATH} -t ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION} -f ${TARGET_PATH}/deploy/Dockerfile ."
//                 sh "docker tag ${DOCKER_IMAGE}:${APP_VERSION} ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION}"
                sh "docker push ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION}"
                sh "docker rmi ${HARBOR_ADDR}/${DOCKER_IMAGE}:${APP_VERSION} -f"
            }

        }
        stage("Deploy") {
            when {
                allOf {
                    expression { env.APP_VERSION != null }
                }
            }
            steps("Deploy to kubernetes") {
                script {
                        sh "export KUBECONFIG=${env.KUBECONFIG}"
                        sh "sed -i 's/VERSION_NUMBER/${APP_VERSION}/g' ${TARGET_PATH}/deploy/k8s-deployment.yml"
                        sh "kubectl --kubeconfig /root/.kube/config apply -f ${TARGET_PATH}/deploy/k8s-deployment.yml --namespace=${K8S_NAMESPACE}"
                }
            }
        }
    }
    post {
    		always {
    			echo 'One way or another, I have finished'
//     			echo sh(returnStdout: true, script: 'env')
    			deleteDir() /* clean up our workspace */
    		}
    		success {
//     			SendDingding("success")
    			echo 'structure success'
    		}
    		failure {
//     			SendDingding("failure")
    			echo 'structure failure'
    		}
       }
}

void SendDingding(res)
{
	// 输入相应的手机号码，在钉钉群指定通知某个人
	tel_num="13008421234"
	//加签 SEC21af6e1dbf64fc40892f9865976266f31b731a897aae6cab4045f3e748bc8c9b
    //https://oapi.dingtalk.com/robot/send?access_token=dc01e09af1e07cd40a926e8e0d0624d78ecb7091b4f04847b87fee2ea45633e2
	// 钉钉机器人的地址
	dingding_url="https://oapi.dingtalk.com/robot/send\\?access_token\\=dc01e09af1e07cd40a926e8e0d0624d78ecb7091b4f04847b87fee2ea45633e2"

    branchName=""
    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
        branchName="理财项目正式环境 tag=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }
    else if (env.GIT_BRANCH ==~ /^release-([0-9])+\.([0-9])+\.([0-9])+.*/){
        branchName="理财项目预生产环境 tag=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }
    else {
        branchName="理财项目开发环境 branch=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }

    // 发送内容
	json_msg=""
	if( res == "success" ) {
		json_msg='{\\"msgtype\\":\\"text\\",\\"text\\":{\\"content\\":\\"@' + tel_num +' [送花花] ' + "${branchName} 第${env.BUILD_NUMBER}次构建，"  + '构建成功。 \\"},\\"at\\":{\\"atMobiles\\":[\\"' + tel_num + '\\"],\\"isAtAll\\":false}}'
	}
	else {
		json_msg='{\\"msgtype\\":\\"text\\",\\"text\\":{\\"content\\":\\"@' + tel_num +' [大哭] ' + "${branchName} 第${env.BUILD_NUMBER}次构建，"  + '构建失败，请及时处理！ \\"},\\"at\\":{\\"atMobiles\\":[\\"' + tel_num + '\\"],\\"isAtAll\\":false}}'
	}

    post_header="Content-Type:application/json;charset=utf-8"
    sh_cmd="curl -X POST " + dingding_url + " -H " + "\'" + post_header + "\'" + " -d " + "\""  + json_msg + "\""
// 	sh sh_cmd
}
