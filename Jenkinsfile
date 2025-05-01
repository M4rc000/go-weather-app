pipeline {
  agent any

  environment {
    GO111MODULE = 'on'
  }

  stages {
    stage('Clone') {
      steps {
        git 'https://your-git-repo.com/project.git'
      }
    }
    stage('Build') {
      steps {
        sh 'go build -o weather-app .'
      }
    }
    stage('Test') {
      steps {
        sh 'go test ./...'
      }
    }
    stage('Docker Build') {
      steps {
        sh 'docker build -t weather-app .'
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploy to production (or staging)...'
      }
    }
  }
}