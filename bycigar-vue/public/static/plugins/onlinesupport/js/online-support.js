// IM BYCMALL PORT START: mall online support ticket and unread badge integration
// EVA START 20260227 新增 支持系统API 客户查询接口
function appendQueryParams(url, params) {
    if (!params) return url;
    var parts = [];
    for (var key in params) {
        if (!Object.prototype.hasOwnProperty.call(params, key)) continue;
        var value = params[key];
        if (value === undefined || value === null || value === '') continue;
        parts.push(encodeURIComponent(key) + '=' + encodeURIComponent(String(value)));
    }
    if (!parts.length) return url;
    return url + (url.indexOf('?') === -1 ? '?' : '&') + parts.join('&');
}
// EVA END 20260227 新增 支持系统API 客户查询接口

// 在线客服聊天窗口初始化
function initOnlineSupport(chatUrl, buttonLogo, buttonText, chatType,nameText) {
    if (chatType == 'omnichat') {
        console.log('omnichat')
        omnichatInit(chatUrl);
    }else if(chatType == 'custom'){
        console.log('custom')
        customChatInit(chatUrl, buttonLogo, buttonText,nameText)
    }else{
        console.log('other')
        other(chatUrl, buttonLogo, buttonText,nameText);
    }
}

function customChatInit(chatUrl, buttonLogo, buttonText, nameText) {
    var supportBtn = document.createElement('div');
    supportBtn.id = 'online-support-btn';
    supportBtn.style.position = 'fixed';
    supportBtn.style.bottom = '20px';
    supportBtn.style.right = '20px';
    supportBtn.style.width = '60px';
    supportBtn.style.height = '60px';
    supportBtn.style.borderRadius = '50%';
    supportBtn.style.backgroundColor = '#4CAF50';
    supportBtn.style.color = 'white';
    supportBtn.style.display = 'flex';
    supportBtn.style.alignItems = 'center';
    supportBtn.style.justifyContent = 'center';
    supportBtn.style.boxShadow = '0 2px 10px rgba(0,0,0,0.2)';
    supportBtn.style.cursor = 'pointer';
    supportBtn.style.zIndex = '9999';
    supportBtn.style.overflow = 'visible';

    if (buttonLogo) {
        supportBtn.style.backgroundImage = 'url("' + buttonLogo + '")';
        supportBtn.style.backgroundSize = 'cover';
        supportBtn.style.backgroundPosition = 'center';
        supportBtn.style.color = 'transparent';
    } else {
        supportBtn.style.fontSize = '24px';
        supportBtn.innerText = buttonText;
    }

    var unreadBadge = document.createElement('span');
    unreadBadge.id = 'online-support-unread';
    unreadBadge.style.position = 'absolute';
    unreadBadge.style.top = '-4px';
    unreadBadge.style.right = '-4px';
    unreadBadge.style.minWidth = '18px';
    unreadBadge.style.height = '18px';
    unreadBadge.style.padding = '0 4px';
    unreadBadge.style.borderRadius = '999px';
    unreadBadge.style.backgroundColor = '#ef4444';
    unreadBadge.style.color = '#fff';
    unreadBadge.style.fontSize = '12px';
    unreadBadge.style.lineHeight = '18px';
    unreadBadge.style.textAlign = 'center';
    unreadBadge.style.display = 'none';
    unreadBadge.style.boxShadow = '0 2px 6px rgba(0,0,0,0.2)';
    supportBtn.appendChild(unreadBadge);

    var chatWindow = null;
    var chatIframe = null;
    var chatWindowVisible = false;
    var pendingOpen = false;
    var expectedOrigin = null;
    try {
        expectedOrigin = new URL(chatUrl, window.location.href).origin;
    } catch (e) {
        expectedOrigin = null;
    }

    function updateUnreadBadge(count) {
        var value = parseInt(count, 10);
        if (!isFinite(value) || value <= 0 || chatWindowVisible) {
            unreadBadge.style.display = 'none';
            unreadBadge.textContent = '';
            return;
        }
        unreadBadge.textContent = value > 99 ? '99+' : String(value);
        unreadBadge.style.display = 'inline-block';
    }

    function postVisibility(open) {
        if (!chatIframe || !chatIframe.contentWindow) return;
        try {
            chatIframe.contentWindow.postMessage({ type: 'support-visibility', open: !!open }, expectedOrigin || '*');
        } catch (e) {
            // ignore
        }
    }

    function hideChatWindow() {
        if (!chatWindow) return;
        chatWindowVisible = false;
        chatWindow.style.opacity = '0';
        chatWindow.style.transform = 'translateY(20px)';
        chatWindow.style.pointerEvents = 'none';
        postVisibility(false);
    }

    function showChatWindow() {
        if (!chatWindow || chatWindowVisible) return;
        chatWindowVisible = true;
        chatWindow.style.pointerEvents = 'auto';
        chatWindow.style.opacity = '1';
        chatWindow.style.transform = 'translateY(0)';
        updateUnreadBadge(0);
        postVisibility(true);
    }

    supportBtn.onclick = function() {
        if (!chatWindow) {
            pendingOpen = true;
            return;
        }
        showChatWindow();
    };

    var supportData = window.initOnlineSupportData || {};

    function openChatWindowWithData(ticketData) {
        if (chatWindow) return;
        var resolvedTicket = ticketData && ticketData.ticket ? ticketData.ticket : (supportData.ticket || '');
        var resolvedCustomerId = ticketData && ticketData.customer_id ? ticketData.customer_id : (supportData.customerId || '');
        var resolvedLogin = ticketData && typeof ticketData.is_login !== 'undefined'
            ? !!ticketData.is_login
            : !!supportData.isLogin;

        function getResponsiveDimensions() {
            var width, height;
            var screenWidth = window.innerWidth;
            var screenHeight = window.innerHeight;

            if (screenWidth < 768) {
                width = Math.min(screenWidth - 40, 320) + 'px';
                height = Math.min(screenHeight * 0.55, 480) + 'px';
            } else {
                width = '380px';
                height = '520px';
            }

            return { width: width, height: height };
        }

        function handleMobileKeyboard() {
            var isMobile = window.innerWidth < 768;
            if (!isMobile) return;

            var originalBottom = '90px';

            function onResize() {
                var currentHeight = window.innerHeight;
                var orientation = window.innerWidth > window.innerHeight ? 'landscape' : 'portrait';
                var viewportHeight = window.screen.height;
                var heightRatio = currentHeight / viewportHeight;

                if (heightRatio < 0.75 && chatWindow.style.bottom !== 'auto' && chatWindow.style.right !== 'auto') {
                    var availableHeight = currentHeight - 120;
                    var windowHeight = parseInt(chatWindow.style.height);

                    if (availableHeight > windowHeight) {
                        chatWindow.style.bottom = (currentHeight - windowHeight - 10) + 'px';
                        chatWindow.style.transition = 'bottom 0.3s ease';
                    }
                } else if (heightRatio > 0.85 && chatWindow.style.bottom !== 'auto' && chatWindow.style.right !== 'auto') {
                    chatWindow.style.bottom = originalBottom;
                    chatWindow.style.transition = 'bottom 0.3s ease';
                }
            }

            window.addEventListener('resize', onResize);

            if (chatIframe) {
                setTimeout(function() {
                    try {
                        var iframeDoc = chatIframe.contentDocument || chatIframe.contentWindow.document;
                        if (iframeDoc) {
                            iframeDoc.addEventListener('focusin', function(e) {
                                if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
                                    setTimeout(onResize, 100);
                                }
                            });
                            iframeDoc.addEventListener('focusout', function(e) {
                                if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
                                    setTimeout(onResize, 300);
                                }
                            });
                        }
                    } catch (err) {
                        console.log('Cannot access iframe content due to cross-origin restrictions');
                    }
                }, 500);
            }

            return function cleanup() {
                window.removeEventListener('resize', onResize);
            };
        }

        chatWindow = document.createElement('div');
        chatWindow.id = 'chat-window';
        chatWindow.style.position = 'fixed';
        chatWindow.style.bottom = '90px';
        chatWindow.style.right = '20px';

        var dimensions = getResponsiveDimensions();
        chatWindow.style.width = dimensions.width;
        chatWindow.style.height = dimensions.height;

        chatWindow.style.backgroundColor = '#1a1a1a';
        chatWindow.style.borderRadius = '12px';
        chatWindow.style.overflow = 'hidden';
        chatWindow.style.boxShadow = '0 10px 40px rgba(0, 0, 0, 0.6)';
        chatWindow.style.zIndex = '9999';
        chatWindow.style.opacity = '0';
        chatWindow.style.transform = 'translateY(20px)';
        chatWindow.style.pointerEvents = 'none';
        chatWindow.style.transition = 'opacity 0.3s ease, transform 0.3s ease';

        var windowHeader = document.createElement('div');
        windowHeader.style.backgroundColor = '#2a2a2a';
        windowHeader.style.padding = '12px 15px';
        windowHeader.style.cursor = 'move';
        windowHeader.style.userSelect = 'none';
        windowHeader.style.display = 'flex';
        windowHeader.style.alignItems = 'center';
        windowHeader.style.justifyContent = 'space-between';
        windowHeader.style.borderBottom = '1px solid #3a3a3a';

        var windowTitle = document.createElement('span');
        windowTitle.textContent = nameText;
        windowTitle.style.color = 'white';
        windowTitle.style.fontSize = '14px';
        windowTitle.style.fontWeight = '500';

        var closeBtn = document.createElement('div');
        closeBtn.style.width = '28px';
        closeBtn.style.height = '28px';
        closeBtn.style.backgroundColor = '#444';
        closeBtn.style.color = 'white';
        closeBtn.style.borderRadius = '50%';
        closeBtn.style.display = 'flex';
        closeBtn.style.alignItems = 'center';
        closeBtn.style.justifyContent = 'center';
        closeBtn.style.fontSize = '16px';
        closeBtn.style.cursor = 'pointer';
        closeBtn.innerHTML = '×';
        closeBtn.style.transition = 'background-color 0.2s ease';

        closeBtn.onmouseenter = function() {
            this.style.backgroundColor = '#666';
        };
        closeBtn.onmouseleave = function() {
            this.style.backgroundColor = '#444';
        };

        var windowContent = document.createElement('div');
        windowContent.style.width = '100%';
        windowContent.style.height = 'calc(100% - 46px)';
        windowContent.style.overflow = 'hidden';

        chatIframe = document.createElement('iframe');

        var finalChatUrl = appendQueryParams(chatUrl, {
            ticket: resolvedTicket || '',
            is_login: resolvedLogin ? 1 : 0,
            ts: Date.now()
        });

        console.log('[OnlineSupport] finalChatUrl:', finalChatUrl);
        chatIframe.src = finalChatUrl;
        chatIframe.style.width = '100%';
        chatIframe.style.height = '100%';
        chatIframe.style.border = 'none';
        chatIframe.style.backgroundColor = '#1a1a1a';
        chatIframe.onload = function() {
            postVisibility(chatWindowVisible);
        };

        windowHeader.appendChild(windowTitle);
        windowHeader.appendChild(closeBtn);
        windowContent.appendChild(chatIframe);
        chatWindow.appendChild(windowHeader);
        chatWindow.appendChild(windowContent);
        document.body.appendChild(chatWindow);

        closeBtn.onclick = function() {
            hideChatWindow();
        };

        var isDragging = false;
        var offsetX, offsetY;

        windowHeader.addEventListener('mousedown', function(e) {
            isDragging = true;
            offsetX = e.clientX - chatWindow.getBoundingClientRect().left;
            offsetY = e.clientY - chatWindow.getBoundingClientRect().top;
            chatWindow.style.transition = 'none';
            chatWindow.style.boxShadow = '0 15px 50px rgba(0, 0, 0, 0.7)';
            chatWindow.style.zIndex = '10001';
            document.body.style.userSelect = 'none';
        });

        document.addEventListener('mousemove', function(e) {
            if (!isDragging) return;

            var newX = e.clientX - offsetX;
            var newY = e.clientY - offsetY;
            var viewportWidth = window.innerWidth;
            var viewportHeight = window.innerHeight;

            newX = Math.max(10, Math.min(newX, viewportWidth - chatWindow.offsetWidth - 10));
            newY = Math.max(10, Math.min(newY, viewportHeight - chatWindow.offsetHeight - 10));

            chatWindow.style.left = newX + 'px';
            chatWindow.style.top = newY + 'px';
            chatWindow.style.bottom = 'auto';
            chatWindow.style.right = 'auto';
        });

        document.addEventListener('mouseup', function() {
            if (isDragging) {
                isDragging = false;
                chatWindow.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
                chatWindow.style.boxShadow = '0 10px 40px rgba(0, 0, 0, 0.6)';
                document.body.style.userSelect = '';
            }
        });

        function handleResize() {
            if (chatWindow.style.bottom !== 'auto' && chatWindow.style.right !== 'auto') {
                var dimensions = getResponsiveDimensions();
                chatWindow.style.width = dimensions.width;
                chatWindow.style.height = dimensions.height;
            }
        }

        window.addEventListener('resize', handleResize);
        handleMobileKeyboard();

        if (pendingOpen) {
            pendingOpen = false;
            showChatWindow();
        }
    }

    if (window.fetch) {
        console.log('发送GET请求获取票据');
        var baseUrl = '';
        if (window.urls && window.urls.base_url) {
            baseUrl = window.urls.base_url;
        } else {
            var path = window.location.pathname || '/';
            var match = path.match(/^\/([a-zA-Z]{2,5}(?:-[a-zA-Z]{2,5})?)\b/);
            var prefix = match ? '/' + match[1] : '';
            baseUrl = window.location.origin + prefix;
        }
        var ticketUrl = baseUrl.replace(/\/$/, '') + '/support/ticket';
        console.log('ticketUrl:', ticketUrl);
        fetch(ticketUrl, { method: 'GET', credentials: 'same-origin' })
            .then(function(res) { return res.json(); })
            .then(function(data) {
                if (data && data.success && data.data) {
                    openChatWindowWithData(data.data);
                } else {
                    openChatWindowWithData(null);
                }
            })
            .catch(function() {
                openChatWindowWithData(null);
            });
    } else {
        openChatWindowWithData(null);
    }

    document.body.appendChild(supportBtn);

    window.addEventListener('message', function(event) {
        if (expectedOrigin && event.origin !== expectedOrigin) return;
        var data = event.data || {};
        if (data.type === 'support-unread') {
            updateUnreadBadge(data.count);
        }
    });
}
// IM BYCMALL PORT END: mall online support ticket and unread badge integration

function omnichatInit(chatUrl) {
    var a=document.createElement('a');
    a.setAttribute('href','javascript:;');
    a.setAttribute('id','easychat-floating-button');
    var span=document.createElement('span');
    span.setAttribute('id', 'easychat-unread-badge');
    span.setAttribute('style','display: none');
    var d1=document.createElement('div');
    d1.setAttribute('id','easychat-close-btn');
    d1.setAttribute('class','easychat-close-btn-close');
    var d2=document.createElement('div');
    d2.setAttribute('id','easychat-chat-dialog');
    d2.setAttribute('class','easychat-chat-dialog-close');
    var ifrm=document.createElement('iframe');
    ifrm.setAttribute('id','easychat-chat-dialog-iframe');
    ifrm.setAttribute('src',chatUrl);
    ifrm.style.width='100%';
    ifrm.style.height='100%';
    ifrm.style.frameborder='0';
    ifrm.style.scrolling='on';
    d2.appendChild(ifrm);
    if(!document.getElementById("easychat-floating-button")){
        document.body.appendChild(a);
        document.body.appendChild(span);
        document.body.appendChild(d1);
        document.body.appendChild(d2);
    }

    var scriptURL = 'https://chat-plugin.easychat.co/easychat.js';
    if(!document.getElementById("omnichat-plugin")) {
    var scriptTag = document.createElement('script');
        scriptTag.src = scriptURL;
        scriptTag.id = 'omnichat-plugin';
        document.body.appendChild(scriptTag);
    }
}

function other(chatUrl, buttonLogo, buttonText,nameText){
    console.log(111)
    // 直接运行传入的 chatUrl 参数作为 JavaScript 代码
    if (chatUrl) {
        try {
            // 检查chatUrl是否包含HTML标签
            if (chatUrl.includes('<script') && chatUrl.includes('</script>')) {
                // 如果包含script标签，提取其中的JavaScript代码
                const scriptContent = chatUrl.replace(/<script[^>]*>/g, '').replace(/<\/script>/g, '');
                // 创建一个脚本元素来执行代码
                var scriptTag = document.createElement('script');
                scriptTag.type = 'text/javascript';
                scriptTag.textContent = scriptContent;
                document.body.appendChild(scriptTag);
                // 执行完后移除脚本元素
                document.body.removeChild(scriptTag);
            } else {
                // 如果不包含script标签，直接执行
                var scriptTag = document.createElement('script');
                scriptTag.type = 'text/javascript';
                scriptTag.textContent = chatUrl;
                document.body.appendChild(scriptTag);
                // 执行完后移除脚本元素
                document.body.removeChild(scriptTag);
            }
        } catch (e) {
            console.error('Error executing chat script:', e);
        }
    }
}

// 在DOM加载完成后初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function() {
        if (window.initOnlineSupportData) {
            initOnlineSupport(
                window.initOnlineSupportData.chatUrl,
                window.initOnlineSupportData.buttonLogo,
                window.initOnlineSupportData.buttonText,
                window.initOnlineSupportData.chatType,
                window.initOnlineSupportData.nameText
            );
        }
    });
} else if (window.initOnlineSupportData) {
    // 如果DOM已经加载完成，直接初始化
    initOnlineSupport(
        window.initOnlineSupportData.chatUrl,
        window.initOnlineSupportData.buttonLogo,
        window.initOnlineSupportData.buttonText,
        window.initOnlineSupportData.chatType,
        window.initOnlineSupportData.nameText
    );
}
