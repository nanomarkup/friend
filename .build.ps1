# Synopsis: Run build
task build {
    if ($PSVersionTable.Platform -ne 'Unix') {
        $AppName = 'friend.exe'
    } else {
        $AppName = 'friend'
    }
    $Status = Start-Process -FilePath 'go' -ArgumentList "build -o ./bin/$AppName" -NoNewWindow -PassThru -Wait
    Assert($Status.ExitCode -eq 0) 'The build command failed'
}

# Synopsis: Run tests
task test {
    $Status = Start-Process -FilePath 'go' -ArgumentList 'test' -NoNewWindow -PassThru -Wait
    Assert($Status.ExitCode -eq 0) 'The test command failed'
}

# Synopsis: Remove generated files
task clean {
    $Status = Start-Process -FilePath 'go' -ArgumentList 'clean' -NoNewWindow -PassThru -Wait 
    Assert($Status.ExitCode -eq 0) 'The "clean" command failed'
}

# Synopsis: Generate documentation
task doc {
    $Status = Start-Process -FilePath 'go' -ArgumentList 'doc -all' -RedirectStandardOutput 'readme.txt' -NoNewWindow -PassThru -Wait 
    Assert($Status.ExitCode -eq 0) 'The "go doc" command failed'
}

task . build, test, clean, doc