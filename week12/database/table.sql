create database test;
use test;

create table if not exists student(
    id int not null auto_increment comment '主键自增id',
    name char(10) not null comment '姓名',
    province char(6) not null comment '省',
    city char(10) not null comment '城市',
    addr varchar(100) default '' comment '地址',
    score float not null default 0 comment '考试成绩',
    enrollment date not null comment '入学时间',
    primary key (id),
    unique key idx_name (name),
    key idx_location (province,city)
)engine=innodb default charset=utf8 comment '学员基本信息';

insert into student (name,province,city,enrollment) values
('张三','北京','北京','2021-03-05'),
('李四','河南','郑州','2021-04-25'),
('小丽','四川','成都','2021-03-10');

select province,avg(score) as avg_score from student where score>0 group by province having avg_score>50 order by avg_score desc;
update student set score=score+10,addr='海淀' where province='北京';
update student set
    score=case province
        when '北京' then score+10 
        when '四川' then score+5 
        else score+7
    end,
    addr=case province
        when '北京' then '东城区'
        when '四川' then '幸福里'
        else '朝阳区'
    end
where id>0;

explain select * from student force index (primary) where province='北京' and id>10000;
