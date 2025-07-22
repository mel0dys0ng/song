#!/bin/bash

# 配置需要操作的目录
DIRECTORIES=("./utils" "../metas" "../erlogs" "../vipers" "../https" "../cobras" "../clients")

# 遍历目录并执行操作
for dir in "${DIRECTORIES[@]}"; do
    echo "正在处理目录: $dir"

    # 切换到目录
    cd "$dir" || { echo "无法切换到目录 $dir"; exit 1; }

    # 输出当前工作目录路径
    echo "当前工作目录: $(pwd)"

    # 执行 git 操作
    echo "执行 git pull..."
    git_pull_output=$(git pull 2>&1)
    echo "$git_pull_output"

    # 执行 go mod tidy 操作
    echo "执行 go mod tidy..."
    go_mod_tidy_output=$(go mod tidy 2>&1)
    echo "$go_mod_tidy_output"

    # 检查 git 是否有未提交的变更
    if [[ -n $(git status --porcelain) ]]; then
        echo "检测到未提交的变更，执行 git 操作..."

        # 执行 git 操作
        echo "执行 git add -A..."
        git_add_output=$(git add -A 2>&1)
        echo "$git_add_output"

        echo "执行 git commit..."
        git_commit_output=$(git commit -am "提交未提交的变更" 2>&1)
        echo "$git_commit_output"

        echo "执行 git push..."
        git_push_output=$(git push 2>&1)
        echo "$git_push_output"

        echo "git 操作完成"
    else
        echo "没有未提交的变更，跳过 git 操作"
    fi

    # 执行 go get -u 并捕获输出
    echo "执行 go get -u..."
    go_get_output=$(go get -u 2>&1)
    echo "$go_get_output"

    # 判断 go get -u 是否有更新
    if [[ -z "$go_get_output" || "$go_get_output" == *"no package to get"* ]]; then
        echo "没有包需要更新，跳过 git 操作"
    else
        echo "检测到包更新，执行 git 操作..."

        # 检查 git 是否有变更
        if [[ -z $(git status --porcelain) ]]; then
            echo "没有变更，跳过 git 操作"
        else
            # 执行 git 操作
            echo "执行 git pull..."
            git_pull_output=$(git pull 2>&1)
            echo "$git_pull_output"

            # 执行 go mod tidy 操作
            echo "执行 go mod tidy..."
            go_mod_tidy_output=$(go mod tidy 2>&1)
            echo "$go_mod_tidy_output"
            
            echo "执行 git add -A..."
            git_add_output=$(git add -A 2>&1)
            echo "$git_add_output"

            echo "执行 git commit..."
            git_commit_output=$(git commit -am "更新依赖包" 2>&1)
            echo "$git_commit_output"

            echo "执行 git push..."
            git_push_output=$(git push 2>&1)
            echo "$git_push_output"

            echo "git 操作完成"
        fi
    fi

    date=$(date 2>&1)
    echo "目录 $dir 处理完成。时间：$date"
    echo "-------------------------"
done

echo "所有目录处理完成"