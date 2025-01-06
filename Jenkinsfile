pipeline {
    agent any
    environment {
        AWS_ACCESS_KEY_ID = credentials('aws-credentials')  // AWS credentials
        AWS_SECRET_ACCESS_KEY = credentials('aws-credentials') // AWS credentials
        INFRACOST_API_KEY = credentials('infracost-api-key') // Infracost API key
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
        stage('Terratest') {
            when {
                expression { env.TF_ACTION == 'apply' }
            }
            steps {
                dir('terraform/tests') {
                       sh '''
                        go mod tidy
                        go test -v
                      '''
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
        stage('Infracost Estimate') {
            when {
                expression { env.TF_ACTION == 'apply' }
            }
            steps {
                dir('terraform') {
                    echo 'Calculating cost estimates using Infracost...'
                    sh '''
                        infracost breakdown --path=. --format=json --out-file=infracost.json
                        infracost output --path=infracost.json --format=table
                    '''
                }
            }
        }
        stage('Terraform Apply') {
            when {
                expression { env.TF_ACTION == 'apply' }
            }
            steps {
                input message: "Proceed with Terraform ${env.TF_ACTION}?", ok: "Continue"
                dir('terraform') {
                    sh 'terraform apply -auto-approve tfplan'
                }
            }
        }
        stage('Terraform Destroy') {
            when {
                expression { env.TF_ACTION == 'destroy' }
            }
            steps {
                input message: "Proceed with Terraform destroy?", ok: "Continue"
                dir('terraform') {
                    sh 'terraform apply -destroy -auto-approve tfplan'
                }
            }
        }
        stage('Visualize Infracost Report') {
            when {
                expression { env.TF_ACTION == 'apply' }
            }
            steps {
                dir('terraform') {
                    echo 'Generating HTML report with cost visualization...'
                    sh '''
                        infracost diff --path=. --format=json --out-file=infracost-diff.json
                        infracost output --path=infracost-diff.json --format=html --out-file=infracost-report.html
                    '''
                }
                archiveArtifacts artifacts: 'terraform/infracost-report.html', allowEmptyArchive: false
                echo "Infracost report has been archived. Review the HTML file for detailed cost breakdown."
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
        // always {
        //     cleanWs(cleanWhenNotBuilt: false, patterns: [[pattern: 'terraform/*.tfstate', type: 'EXCLUDE']])
        // }
    }
}