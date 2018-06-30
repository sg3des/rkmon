	var ajax = {
			init: function () {
				var xmlhttp;
				try {
					xmlhttp = new ActiveXObject("Msxml2.XMLHTTP");
				} catch (e) {
					try {
						xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
					} catch (E) {
						xmlhttp = false;
					}
				}
				if (!xmlhttp && typeof XMLHttpRequest != 'undefined') xmlhttp = new XMLHttpRequest();
				return xmlhttp;
			},

			serialize: function(obj) {
			  var str = [];
			  for (var p in obj)
			    if (obj.hasOwnProperty(p)) {
			      str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
			    }
			  return str.join("&");
			},

			send: function(method, url, data, callback) {
				var a = ajax.init();

				method = method.toUpperCase();
				if (method == "POST" || method == "PUT" || method == "UPDATE") {
					a.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
				} else {
					url += "?"+ajax.serialize(data);
				}

				a.open(method, url);
				a.send(data);
				
				if (callback) {
					a.onreadystatechange = function () {
						if (a.readyState >= 4) {
							if (callback == undefined || callback == "") {
								console.log(a.responseText);
								return;
							}

							switch (callback.charAt(0)) {
								case "#":
									//insert response by id
									var dst = document.getElementById(callback.substring(1, callback.length));
									dst.innerHTML = a.responseText;
									break;

								case ".":
									//insert response by classname
									var dsts = document.getElementsByClassName(callback.substring(1, callback.length));
									dsts.forEach(function(dst) {
									  dst.innerHTML = a.responseText;
									});
									break;

								default:
									//pass response to callback function
									callback(a.responseText);
									break;
							}
							
						}
					}
				}
			},

			serialize: function(obj) {
				var str = [];
				for(var p in obj)
					if (obj.hasOwnProperty(p)) {
						str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
					}
				return str.join("&");
			}
		}


		function checkIP(ip, dst) {
			console.log(ip);
			console.log(dst);
			ajax.send("GET", "/", {"ip": ip}, dst);
		}