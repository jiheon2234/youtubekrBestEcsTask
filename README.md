https://github.com/jiheon2234/youtubeKrBest

유튜브 API를 통해 인기있는 유튜브 영상 정보와 그 댓글들을 가져오고, 이를 db에 적재

ECS 태스크로 실행할 작업
aws 이벤트브리지 cron을 통해서 하루에 1번씩 
- 실패시 alarm+sns로 이메일
- 성공시 보고서만드는 lambda 실행됨

메인에 push하면 github action으로  docker hub에 푸시됨