node {
    checkout scm
    def userId = sh(script: "id -u ${USER}", returnStdout: true).trim()
    def customImage = docker.build("build-go-sudoku:${env.BUILD_ID}", "--build-arg JENKINS_UID=${userId} -f .jenkins/docker/Dockerfile .jenkins/docker")

    sh("""
       if [ -d ${WORKSPACE_TMP}/${env.BUILD_ID}/go ]; then
         find ${WORKSPACE_TMP}/${env.BUILD_ID}/go -type d ! -writable -exec chmod u+w {} ";"
         rm -rf ${WORKSPACE_TMP}/${env.BUILD_ID}/go
       fi
       mkdir -p ${WORKSPACE_TMP}/${env.BUILD_ID}/go
       """)

    withCredentials([string(credentialsId: 'coverage-token', variable: 'COVERAGE_TOKEN')]) {
        cache(maxCacheSize: 250, defaultBranch: 'master', caches: [
            [$class: 'ArbitraryFileCache', path: "${env.WORKSPACE_TMP}/${env.BUILD_ID}/go", cacheValidityDecidingFile: 'go.mod', compressionMethod: 'TARGZ']
        ]) {
            customImage.inside("-v ${env.WORKSPACE_TMP}/${env.BUILD_ID}/go:/home/jenkins/go") {
                stage('Build') {
                   sh 'go build'
                }
                stage('Test') {
                   sh 'go test -coverprofile=coverage.txt -covermode=atomic'
                }
                stage('Upload coverage to codecov') {
                   sh './scripts/codecov.sh -t $COVERAGE_TOKEN -K'
                }
            }
        }
    }
    sh("""
       if [ -d ${WORKSPACE_TMP}/${env.BUILD_ID} ]; then
         find ${WORKSPACE_TMP}/${env.BUILD_ID} -type d ! -writable -exec chmod u+w {} ";"
         rm -rf ${WORKSPACE_TMP}/${env.BUILD_ID}
       fi
       """)
}
