node {
    checkout scm
    def customImage = docker.build("build-go-sudoku:${env.BUILD_ID}", "-f .jenkins/docker/Dockerfile .jenkins/docker")
    withCredentials([string(credentialsId: 'coverage-token', variable: 'COVERAGE_TOKEN')]) {
        customImage.inside('-v $HOME/go:/home/jenkins/go') {
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
