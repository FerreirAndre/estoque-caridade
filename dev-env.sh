#!/bin/bash

session="dev"

tmux new-session -d -s $session

tmux rename-window -t 1 'nvim'
tmux send-keys -t 'nvim' 'nvim .' C-m

tmux new-window -t $session:2 -n 'server-back'
#tmux send-keys -t 'server-back' 'cd ecom-api' C-m 'go run cmd/api/main.go' C-m
#tmux new-window -t $session:3 -n 'server-front'
#tmux send-keys -t 'server-front' 'cd ecom-ui' C-m 'ng s' C-m
tmux new-window -t $session:3 -n 'terminal'

tmux attach-session -t $session:1
