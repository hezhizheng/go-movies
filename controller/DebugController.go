package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go_movies/utils/spider"
	"net/http"
	"regexp"
	"strings"
)

func Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	t := spider.FormatVodPDownUrl("第01集$http://okxzy.xzokzyzy.com/20200406/14933_31886b53/高玩救未来.Future.Man.S03E01.mp4#第02集$http://okxzy.xzokzyzy.com/20200406/14934_87bdd380/高玩救未来.Future.Man.S03E03.mp4#第03集$http://okxzy.xzokzyzy.com/20200406/14935_315c0332/高玩救未来.Future.Man.S03E02.mp4")

	str := "第01集$http://okxzy.xzokzyzy.com/20200406/14933_31886b53/高玩救未来.Future.Man.S03E01.mp4#第02集$http://okxzy.xzokzyzy.com/20200406/14934_87bdd380/高玩救未来.Future.Man.S03E03.mp4#第03集$http://okxzy.xzokzyzy.com/20200406/14935_315c0332/高玩救未来.Future.Man.S03E02.mp4"

	c := strings.Split(str, "$")

	matched, err := regexp.Compile("(((ht|f)tps?):\\/\\/)?[\\w-]+(\\.[\\w-]+)+([\\w.,@?^=%&:/~]*[\\w@?^=%&/~])?")
	fmt.Println(matched.FindAllString(str,-1), err)

	fmt.Fprint(w, "DEBUG", t,c,c[1:])

	return
	//film := "6,7,8,9,10,11,12,20,21,31"
	//t := strings.Contains(film, "21")
	//fmt.Fprint(w, "DEBUG", t)
	//return
	//go utils.DelAllListCacheKey()

	//c := utils.RedisDB.Exists("detail_links:id:10").Val()
	//fmt.Fprint(w, "DEBUG", c)
	//return
	go spider.StartApi()

	return

	//x1 := `HD高清$https://youku.cdn10-okzy.com/share/aec3d6abde16f9aad1970fad0fed2cb3$$$HD高清$https://youku.cdn10-okzy.com/20200406/13525_b65f0b6c/index.m3u8`
	x1 := `第01集$https://iqiyi.cdn9-okzy.com/share/1d01bd2e16f57892f0954902899f0692#第02集$https://iqiyi.cdn9-okzy.com/share/664c7298d2b73b3c7fe2d1e8d1781c06#第03集$https://iqiyi.cdn9-okzy.com/share/6e2290dbf1e11f39d246e7ce5ac50a1e$$$第01集$https://iqiyi.cdn9-okzy.com/20200406/8392_bc901e4e/index.m3u8#第02集$https://iqiyi.cdn9-okzy.com/20200406/8394_d3645c2c/index.m3u8#第03集$https://iqiyi.cdn9-okzy.com/20200406/8393_f1911261/index.m3u8`
	y, y1 := spider.FormatVodPlayUrl(x1)
	fmt.Fprint(w, "DEBUG", y, y1)
	return


}
