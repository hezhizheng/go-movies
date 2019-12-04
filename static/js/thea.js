var browser={    
		versions:function(){            
				var u = navigator.userAgent, app = navigator.appVersion;            
				return {                
					trident: u.indexOf('Trident') > -1,               
					presto: u.indexOf('Presto') > -1,                
					webKit: u.indexOf('AppleWebKit') > -1,              
					gecko: u.indexOf('Gecko') > -1 && u.indexOf('KHTML') == -1,               
					mobile: !!u.match(/AppleWebKit.*Mobile.*/)||!!u.match(/AppleWebKit/),          
					ios: !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/),                 
					android: u.toLowerCase().indexOf('android') > -1 ,   
					iPhone: u.indexOf('iPhone') > -1 || u.indexOf('Mac') > -1,               
					iPad: u.indexOf('iPad') > -1,               
					webApp: u.indexOf('Safari') == -1           
				};
				}()
}
if (!(browser.versions.android || browser.versions.ios || browser.versions.iPhone || browser.versions.iPad)){document.write('<script type=\'text/javascript\' src=\'https://6vhao.kkcaicai.com/960X90.js\' charset=\'GBK\'></script><script type=\'text/javascript\' src=\'https://6vhao.kkcaicai.com/tan.js\' charset=\'GBK\'></script>');
}