# 公主连接装备查看器

> [访问地址](https://pcr-go-server.us-south.cf.appdomain.cloud/)
>
> 前端采用vue编写
>
> 项目采用IBM Cloud Foundry进行项目的自动化部署 

## 实现功能

- 通过点选装备，输出装备的掉落列表
    - 根据装备的爆率以及同一张地图出现所需装备的多寡排序掉落地图，HARD地图默认优先级最低
    - 可以对点选的装备进行排序，越前面的装备权重越高
- 角色Rank推荐
