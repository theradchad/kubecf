#!/usr/bin/env ruby

# Doc ref: https://docs.bazel.build/versions/master/user-manual.html#flag--workspace_status_command.

git_root = `git rev-parse --show-toplevel`.strip!

# Embed the workspace status from the credhub_setup workspace.
require "#{git_root}/deploy/containers/credhub_setup/workspace_status.rb"
