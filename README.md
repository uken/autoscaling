- Create docker image for app
    docker run -a stdout -i -v $PWD/Bazingo:/build -v $PWD/cache:/cache lxfontes/frontend_19:0.0.1 run /build/ops/build.sh
- Run seed / migrations
    docker run -a stdout -i -v $PWD/Bazingo:/build -v $PWD/cache:/cache lxfontes/frontend_19:0.0.1 run bundle exec rake db:migrate db:seed
- Start
    docker run -a stdout -i -v $PWD/Bazingo:/build -v $PWD/cache:/cache lxfontes/frontend_19:0.0.1 start -c all=0,web=1

