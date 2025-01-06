pipeline {
    agent any
    environment {
        AWS_ACCESS_KEY_ID = credentials('aws-credentials')  // Replace with your Jenkins credential ID
        AWS_SECRET_ACCESS_KEY = credentials('aws-credentials') // Replace with your Jenkins credential ID
    }
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        stage('Terraform Init') {
            steps {
                dir('terraform') {
                    sh 'terraform init'
                }
            }
        }
        stage('Terraform Validate') {
            steps {
                dir('terraform') {
                    sh 'terraform validate'
                }
            }
        }
        stage('Select Action') {
            steps {
                script {
                    env.TF_ACTION = input(
                        id: 'ActionChoice', message: 'Select Terraform Action:',
                        parameters: [choice(name: 'ACTION', choices: 'apply\ndestroy', description: 'Choose whether to apply or destroy infrastructure')]
                    )
                }
            }
        }
        stage('Terraform Plan') {
            steps {
                dir('terraform') {
                    script {
                        if (env.TF_ACTION == 'apply') {
                            sh 'terraform plan -out=tfplan'
                        } else {
                            sh 'terraform plan -destroy -out=tfplan'
                        }
                    }
                }
            }
        }
        stage('Terraform Apply/Destroy') {
            steps {
                input message: "Proceed with Terraform ${env.TF_ACTION}?", ok: "Continue"
                dir('terraform') {
                    script {
                        if (env.TF_ACTION == 'apply') {
                            sh 'terraform apply -auto-approve tfplan'
                        } else {
                            sh 'terraform apply -destroy -auto-approve tfplan'
                        }
                    }
                }
            }
        }
    }
    post {
        success {
            echo "Terraform ${env.TF_ACTION} completed successfully."
        }
        failure {
            echo "Terraform ${env.TF_ACTION} failed. Check logs for details."
        }
    }
}
