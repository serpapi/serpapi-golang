# serpapi-golang release package

require 'fileutils'

def version
  @version ||= begin
    version_line = File.readlines('serpapi.go').find { |line| line.include?('VERSION') && line.include?('=') }
    if version_line
      match = version_line.match(/VERSION\s*=\s*"([^"]+)"/)
      match[1] if match
    end
  end
end

desc "Run all tasks: version, lint, test, doc, ready"
task :default => [:version, :lint, :test, :doc, :ready]

desc "Lint source code using go tools"
task :lint => [:vet, :format]

desc "Run go vet"
task :vet do
  puts "Run go vet"
  sh "go vet ."
  Dir.glob('./demo/*').each do |file|
    sh "go vet #{file}"
  end
  sh "go vet ./test"
end

desc "Format code"
task :format do
  sh "go fmt ."
  sh "go fmt ./test"
  Dir.glob('./demo/*').each do |file|
    sh "go fmt #{file}"
  end
end

desc "Run integration test suite"
task :test do
  sh "go test -v ./test"
end

desc "Run code coverage"
task :coverage do
  puts "Run code coverage"
  sh "go test -cover -covermode=count -coverpkg=./... -coverprofile=coverage.out ./test"
  sh "go tool cover -html=coverage.out -o coverage.html"
  puts "Coverage report generated: coverage.html"
end

desc "Run examples"
task :example do
  puts "run example"
  sh "go test -v ./test/example/*.go"
end

desc "Create documentation"
task :doc do
  sh "go doc"
end

desc "Check that everything is pushed"
task :ready do
  puts "check if repository has changes"
  result = `git status`
  unless result.downcase.include?("nothing") || result.include?("working tree clean")
    raise "Repository has uncommitted changes"
  end
end

desc "Out of box testing - validate the pre-released library"
task :oobt do
  FileUtils.mkdir_p '/tmp/serpapi-golang'
  FileUtils.cp 'oobt/demo.go', '/tmp/serpapi-golang' if File.exist?('oobt/demo.go')
  Dir.chdir('/tmp/serpapi-golang') do
    sh "go mod init serpapi.com/golang/oobt"
    sh "go get -u github.com/serpapi/serpapi-golang"
    sh "go run demo.go"
  end
end

desc "Show current version for golang and library"
task :version do
  puts "golang: #{`go version`.strip}"
  v = version || "unknown"
  puts "current version: #{v}"
end

desc "Display the current release information"
task :release => [:oobt, :version] do
  v = version
  raise "Could not determine version" unless v
  # tags are prefixed with v (vX.Y.Z) as required by the go module resolver
  tag = "v#{v}"
  sh "git tag -a #{tag}"
  sh "git push origin #{tag}"
  puts "create release: #{tag}"
end

desc "Clean build artifacts"
task :clean do
  FileUtils.rm_f 'coverage.out'
  FileUtils.rm_f 'coverage.html'
  sh "go clean"
end
