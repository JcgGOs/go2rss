package util

import (
	"fmt"
	"regexp"
	"testing"
)

func TestClearComment(t *testing.T) {

	tests := []struct {
		name string
		args string
		want string
	}{
		{"1", "<!-- xx -->", ""},
		{"2", "<p><!-- xx --></p>", "<p></p>"},
		{"3", "<p><!-- xx --><!-- xx --><!-- xx --><!-- xx --></p>", "<p></p>"},
		{"4", "<p><!-- xx ---></p>", "<p></p>"},
		{"5", "<p><!-- xx ---></p><p><!-- xx ---></p>", "<p></p><p></p>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClearComment(tt.args); got != tt.want {
				t.Errorf("ClearComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

var reTmp = regexp.MustCompile(`<!--[\s\S]*?-->`)
var reTmp2 = regexp.MustCompile(`<script>[\s\S]*?</script>`)

func TestClearComment2(t *testing.T) {

	str := `
	<!--top
	--><!--content--><p>　　<strong>2022成都
	地铁红包墙活动时间</strong></p><p>　　成都地铁联合支付宝推出线上
	乘车10元礼券，参与以下四场线下活动，通过现场扫码即可领取</p><p>　　<strong>活动时
	间及地点</strong></p><p>　　<strong>1月20日、27日在世纪城</strong></p><p>　　<strong>1月2
	1日、28日在火车南站当日10:30—16:30</strong></p><p>　　车站将设置有“红包墙”，乘客扫描红包背后的二维码即可获得线上乘车10元礼券</p><p>　　该券自领取之时起7×24小时内有效</p><p>　　可在成都地铁APP使用支付宝代扣或使用支付宝APP乘车码乘坐地铁时直接抵扣</p><p>　　不限使用次数，且单笔订单不限额，金额用完或券到期为止</p><p>　　详细规则请见支付宝领取界面或详询支付宝客服</p><div id="adInArticle"></div><div class="ad-in-article a-2"><!-- <div style="width:760px;height:100px;float:left;margin: 14px 0;"> <script type="text/javascript" src="http://d.s11.cn/x7dry1fhrn.js"
	></script> </div> ---></div><!--bottom--><div id="add_ewm_content"><p>温馨提示：微信搜索公众号成都本地宝，关注后在对话框回复【地铁】，即可查询成都地铁建设进展、在建线路规划、地铁运营时间、轨道交通线路等消息。</p><p style="text-align:center;"><img src="http://imgbdb3.bendibao.com/dazheimg/202111/12/20211112105842_67935.png" width="244" height="400"/></p></div><script>var newNode = document.createElement("div");newNode.innerHTML =document.getElementById("add_ewm_content").innerHTML ;document.getElementById("add_ewm_content").remove();document.getElementById("adInArticle").insertBefore(newNode,document.getElementById("adInArticle").firstChild);</script><!--mobile--><p class="view_city_index"><a style="background:none;" href="http://m.cd.bendibao.com/traffic/133308.shtm" target="_blank">手机访问</a> <a href="http://cd.bendibao.com" target="_blank">成都本地宝首页</a></p><div id="sclear" style="clear:both;"></div>	
	`
	context := reTmp.ReplaceAllString(str, "")
	fmt.Printf("context: %v\n", context)

	context = reTmp2.ReplaceAllString(context, "")
	fmt.Printf("context: %v\n", context)
}
