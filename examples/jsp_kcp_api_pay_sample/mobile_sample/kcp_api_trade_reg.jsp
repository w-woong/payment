<%@ page language="java" contentType="text/html;charset=euc-kr"%>
<%@ page import="org.json.simple.JSONObject" %>
<%@ page import="org.json.simple.parser.JSONParser" %>
<%@ page import="org.json.simple.parser.ParseException" %>
<%@ page import="java.io.*" %>
<%@ page import="java.util.*" %>
<%@ page import="java.net.HttpURLConnection" %>
<%@ page import="java.net.URL" %>
<%!
    /* 
    ==========================================================================
    null 값을 처리하는 메소드                                                               
    --------------------------------------------------------------------------
    */
        public String f_get_parm( String val )
        {
          if ( val == null ) val = "";
          return  val;
        }
    /* ============================================================================== */
%>
<%
    request.setCharacterEncoding( "euc-kr" );
    /* 
    ==========================================================================
          거래등록 API URL
    --------------------------------------------------------------------------
    */
    String target_URL = "https://stg-spl.kcp.co.kr/std/tradeReg/register"; //개발환경
    //String target_URL = "https://spl.kcp.co.kr/std/tradeReg/register"; //운영환경
    /* 
    ==========================================================================
          요청 정보                                                          
    --------------------------------------------------------------------------
    */
    String site_cd              = f_get_parm( request.getParameter( "site_cd"  )); // 사이트코드
    // 인증서정보(직렬화)
    String kcp_cert_info        = "-----BEGIN CERTIFICATE-----MIIDgTCCAmmgAwIBAgIHBy4lYNG7ojANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJLUjEOMAwGA1UECAwFU2VvdWwxEDAOBgNVBAcMB0d1cm8tZ3UxFTATBgNVBAoMDE5ITktDUCBDb3JwLjETMBEGA1UECwwKSVQgQ2VudGVyLjEWMBQGA1UEAwwNc3BsLmtjcC5jby5rcjAeFw0yMTA2MjkwMDM0MzdaFw0yNjA2MjgwMDM0MzdaMHAxCzAJBgNVBAYTAktSMQ4wDAYDVQQIDAVTZW91bDEQMA4GA1UEBwwHR3Vyby1ndTERMA8GA1UECgwITG9jYWxXZWIxETAPBgNVBAsMCERFVlBHV0VCMRkwFwYDVQQDDBAyMDIxMDYyOTEwMDAwMDI0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAppkVQkU4SwNTYbIUaNDVhu2w1uvG4qip0U7h9n90cLfKymIRKDiebLhLIVFctuhTmgY7tkE7yQTNkD+jXHYufQ/qj06ukwf1BtqUVru9mqa7ysU298B6l9v0Fv8h3ztTYvfHEBmpB6AoZDBChMEua7Or/L3C2vYtU/6lWLjBT1xwXVLvNN/7XpQokuWq0rnjSRThcXrDpWMbqYYUt/CL7YHosfBazAXLoN5JvTd1O9C3FPxLxwcIAI9H8SbWIQKhap7JeA/IUP1Vk4K/o3Yiytl6Aqh3U1egHfEdWNqwpaiHPuM/jsDkVzuS9FV4RCdcBEsRPnAWHz10w8CX7e7zdwIDAQABox0wGzAOBgNVHQ8BAf8EBAMCB4AwCQYDVR0TBAIwADANBgkqhkiG9w0BAQsFAAOCAQEAg9lYy+dM/8Dnz4COc+XIjEwr4FeC9ExnWaaxH6GlWjJbB94O2L26arrjT2hGl9jUzwd+BdvTGdNCpEjOz3KEq8yJhcu5mFxMskLnHNo1lg5qtydIID6eSgew3vm6d7b3O6pYd+NHdHQsuMw5S5z1m+0TbBQkb6A9RKE1md5/Yw+NymDy+c4NaKsbxepw+HtSOnma/R7TErQ/8qVioIthEpwbqyjgIoGzgOdEFsF9mfkt/5k6rR0WX8xzcro5XSB3T+oecMS54j0+nHyoS96/llRLqFDBUfWn5Cay7pJNWXCnw4jIiBsTBa3q95RVRyMEcDgPwugMXPXGBwNoMOOpuQ==-----END CERTIFICATE-----";
    String ordr_idxx            = f_get_parm( request.getParameter( "ordr_idxx" )); // 주문번호
    String good_mny             = f_get_parm( request.getParameter( "good_mny" )); // 결제 금액
    String good_name            = f_get_parm( request.getParameter( "good_name" )); // 상품명
    String pay_method           = f_get_parm( request.getParameter( "pay_method" )); // 결제수단
    String Ret_URL              = f_get_parm( request.getParameter( "Ret_URL" )); // 리턴 URL
    /* ============================================================================== */
    String actionResult         = f_get_parm( request.getParameter( "ActionResult" )); // pay_method에 매칭되는 값 (인증창 호출 시 필요)
    String van_code             = f_get_parm( request.getParameter( "van_code" )); // (포인트,상품권 인증창 호출 시 필요)
    
    JSONObject json_req = new JSONObject();

    json_req.put("site_cd", site_cd);
    json_req.put("kcp_cert_info", kcp_cert_info);
    json_req.put("ordr_idxx", ordr_idxx);
    json_req.put("good_mny", good_mny);
    json_req.put("good_name", good_name);
    json_req.put("pay_method", pay_method);
    json_req.put("Ret_URL", Ret_URL);
    json_req.put("escw_used", "N");
    json_req.put("user_agent", "");

    String temp_req_param = json_req.toString();
    String req_param = temp_req_param.replace(",",",\r\n");

    String inputLine = null;
    StringBuffer outResult = new StringBuffer();

    try 
    {
        // API REQ
        URL url = new URL(target_URL);
        HttpURLConnection conn = (HttpURLConnection)url.openConnection();
        conn.setDoOutput(true);
        conn.setRequestMethod("POST");
        conn.setRequestProperty("Content-Type", "application/json");
        conn.setRequestProperty("Accept-Charset", "UTF-8");

        OutputStream os = conn.getOutputStream();
        os.write(req_param.getBytes("UTF-8"));
        os.flush();
        
        // API RES
        BufferedReader in = new BufferedReader(new InputStreamReader(conn.getInputStream(), "UTF-8"));
        while ((inputLine = in.readLine()) != null) 
        {
            outResult.append(inputLine);
        }
        conn.disconnect();
    }
    catch(Exception e)
    {
        e.printStackTrace();
    }

    String temp_result = outResult.toString();
    String res_param = temp_result.replace(",",",\r\n");
    
    /* 
    ==========================================================================
    거래등록 응답정보                                                               
    --------------------------------------------------------------------------
    */
    String res_cd      = ""; // 응답코드
    String res_msg     = ""; // 응답메세지
    String approvalKey = ""; // 거래등록키
    String traceNo     = ""; // 추적번호
    String PayUrl      = ""; // 거래등록 PAY URL
    
    // RES JSON DATA Parsing
    JSONParser parser = new JSONParser();
    JSONObject json_res = (JSONObject)parser.parse(temp_result);
    
    res_cd      = (String)json_res.get("Code");
    res_msg     = (String)json_res.get("Message");
    approvalKey = (String)json_res.get("approvalKey");
    traceNo     = (String)json_res.get("traceNo");
    PayUrl      = (String)json_res.get("PayUrl");

%>

<!DOCTYPE>
<html>
<head>
    <title>*** NHN KCP API SAMPLE ***</title>
    <meta http-equiv="Content-Type" content="text/html; charset=euc-kr" />
    <meta http-equiv="x-ua-compatible" content="ie=edge"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=yes, target-densitydpi=medium-dpi">  
    <script type="text/javascript">

    function goReq()
    {
        <%
        // 거래등록 처리 정상
        if ( res_cd.equals( "0000" ) )
        {
        %>
            alert("거래등록 성공");
            document.form_trade_reg.action = "order_mobile.jsp";
            document.form_trade_reg.submit();
        <%
        }
    
        // 거래등록 처리 실패, 여기(샘플)에서는 trade_reg page로 리턴 합니다.
        else
        {
        %>
            alert("에러 코드 : <%=res_cd%>, 에러 메세지 : <%=res_msg%>");
            location.href = "./trade_reg.html";
        <%
        }
        %>
    }
    </script>
</head>
<body onload="goReq();">
    <div class="wrap">
        <!--  거래등록 form : form_trade_reg -->
        <form name="form_trade_reg" method="post">
            <input type="hidden" name="site_cd"         value="<%=site_cd %>" />  <!-- 사이트 코드 -->
            <input type="hidden" name="ordr_idxx"       value="<%=ordr_idxx %>" /><!-- 주문번호     -->
            <input type="hidden" name="good_mny"        value="<%=good_mny %>" /> <!-- 결제금액     -->
            <input type="hidden" name="good_name"       value="<%=good_name %>" /><!-- 상품명        -->
            <!-- 인증시 필요한 파라미터(변경불가)-->
            <input type="hidden" name="pay_method"      value="<%=pay_method %>" />
            <input type="hidden" name="ActionResult"    value="<%=actionResult %>" />
            <input type="hidden" name="van_code"        value="<%=van_code %>" />
            <!-- 리턴 URL (kcp와 통신후 결제를 요청할 수 있는 암호화 데이터를 전송 받을 가맹점의 주문페이지 URL) -->
            <input type="hidden" name="Ret_URL"         value="<%=Ret_URL %>" />
            <!-- 거래등록 응답 값 -->
            <input type="hidden" name="approvalKey"     value="<%=approvalKey %>" />
            <input type="hidden" name="traceNo"         value="<%=traceNo %>" />
            <input type="hidden" name="PayUrl"          value="<%=PayUrl %>" />
        </form>
    </div>
<!--//wrap-->
</body>
</html>