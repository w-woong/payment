# python 샘플의 경우 flask 프레임워크로 구성되어있습니다
from lib2to3.pgen2.pgen import DFAState
from pickle import FALSE
from flask import Flask, render_template, request
import requests, json
import OpenSSL
from OpenSSL import crypto
import base64

app = Flask(__name__)


# 테스트용 인증서정보(직렬화)
KCP_CERT_INFO = '-----BEGIN CERTIFICATE-----MIIDgTCCAmmgAwIBAgIHBy4lYNG7ojANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJLUjEOMAwGA1UECAwFU2VvdWwxEDAOBgNVBAcMB0d1cm8tZ3UxFTATBgNVBAoMDE5ITktDUCBDb3JwLjETMBEGA1UECwwKSVQgQ2VudGVyLjEWMBQGA1UEAwwNc3BsLmtjcC5jby5rcjAeFw0yMTA2MjkwMDM0MzdaFw0yNjA2MjgwMDM0MzdaMHAxCzAJBgNVBAYTAktSMQ4wDAYDVQQIDAVTZW91bDEQMA4GA1UEBwwHR3Vyby1ndTERMA8GA1UECgwITG9jYWxXZWIxETAPBgNVBAsMCERFVlBHV0VCMRkwFwYDVQQDDBAyMDIxMDYyOTEwMDAwMDI0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAppkVQkU4SwNTYbIUaNDVhu2w1uvG4qip0U7h9n90cLfKymIRKDiebLhLIVFctuhTmgY7tkE7yQTNkD+jXHYufQ/qj06ukwf1BtqUVru9mqa7ysU298B6l9v0Fv8h3ztTYvfHEBmpB6AoZDBChMEua7Or/L3C2vYtU/6lWLjBT1xwXVLvNN/7XpQokuWq0rnjSRThcXrDpWMbqYYUt/CL7YHosfBazAXLoN5JvTd1O9C3FPxLxwcIAI9H8SbWIQKhap7JeA/IUP1Vk4K/o3Yiytl6Aqh3U1egHfEdWNqwpaiHPuM/jsDkVzuS9FV4RCdcBEsRPnAWHz10w8CX7e7zdwIDAQABox0wGzAOBgNVHQ8BAf8EBAMCB4AwCQYDVR0TBAIwADANBgkqhkiG9w0BAQsFAAOCAQEAg9lYy+dM/8Dnz4COc+XIjEwr4FeC9ExnWaaxH6GlWjJbB94O2L26arrjT2hGl9jUzwd+BdvTGdNCpEjOz3KEq8yJhcu5mFxMskLnHNo1lg5qtydIID6eSgew3vm6d7b3O6pYd+NHdHQsuMw5S5z1m+0TbBQkb6A9RKE1md5/Yw+NymDy+c4NaKsbxepw+HtSOnma/R7TErQ/8qVioIthEpwbqyjgIoGzgOdEFsF9mfkt/5k6rR0WX8xzcro5XSB3T+oecMS54j0+nHyoS96/llRLqFDBUfWn5Cay7pJNWXCnw4jIiBsTBa3q95RVRyMEcDgPwugMXPXGBwNoMOOpuQ==-----END CERTIFICATE-----'

# INDEX PAGE
@app.route('/')
def index():
    return render_template('index.html')

# ORDER PAGE(PC)
@app.route('/sample/order')
def order():
    return render_template('sample/order.html')

# MOBILE 거래등록 PAGE
@app.route('/mobile_sample/trade_reg')
def trade_reg():
    return render_template('/mobile_sample/trade_reg.html')

# MOBILE 거래등록 API
@app.route('/mobile_sample/kcp_api_trade_reg', methods=['POST'])
def kcp_api_trade_reg():
    #거래등록처리 POST DATA
    actionResult = f_get_parm(request.form['ActionResult']) # pay_method에 매칭되는 값 (인증창 호출 시 필요)
    van_code = f_get_parm(request.form['van_code']) # (포인트,상품권 인증창 호출 시 필요)
    
    post_data = {
        'actionResult' : actionResult, 
        'van_code': van_code
    }
    # 거래등록 API
    target_URL = 'https://stg-spl.kcp.co.kr/std/tradeReg/register' #개발환경
    #target_URL = 'https://spl.kcp.co.kr/std/tradeReg/register' #운영환경
    headers = {'Content-Type': 'application/json', 'charset': 'UTF-8'}

    # 거래등록 API REQ DATA
    req_data = {
        'site_cd' : f_get_parm(request.form['site_cd']),
        'kcp_cert_info' : KCP_CERT_INFO,
        'ordr_idxx' : f_get_parm(request.form['ordr_idxx']),
        'good_mny' : f_get_parm(request.form['good_mny']),
        'good_name' :f_get_parm(request.form['good_name']),
        'pay_method' : f_get_parm(request.form['pay_method']),
        'Ret_URL' : f_get_parm(request.form['Ret_URL']),
        'escw_used' : 'N',
        'user_agent' : ''
    }
    
    res = requests.post(target_URL, headers=headers, data=json.dumps(req_data, ensure_ascii=False, indent="\t").encode('utf8')) 
    return render_template('mobile_sample/kcp_api_trade_reg.html', res_data=json.loads(res.text), req_data=req_data, post_data=post_data)

#주문페이지 이동 및 Ret_URL 처리(MOBILE)
@app.route('/mobile_sample/order_mobile', methods=['POST'])
def order_mobile():
    # res_cd == '0000' 인 경우 진행
    if f_get_parm(request.form['res_cd']) == '0000' : 
        enc_info = f_get_parm(request.form['enc_info'])
        post_data = ''
        # enc_info 값이 없을 경우 POST DATA 처리 후 order_mobile 이동
        if enc_info == '' :
            post_data = {
                'approvalKey' : f_get_parm(request.form['approvalKey']),
                'traceNo' : f_get_parm(request.form['traceNo']),
                'PayUrl' : f_get_parm(request.form['PayUrl']),
                'pay_method' : f_get_parm(request.form['pay_method']),
                'actionResult' : f_get_parm(request.form['ActionResult']),
                'Ret_URL' : f_get_parm(request.form['Ret_URL']),
                'van_code' : f_get_parm(request.form['van_code']),
                'site_cd' : f_get_parm(request.form['site_cd']),
                'ordr_idxx' : f_get_parm(request.form['ordr_idxx']), # 쇼핑몰 주문번호   
                'good_name' : f_get_parm(request.form['good_name']), # 상품명            
                'good_mny' : f_get_parm(request.form['good_mny']) # 결제 금액      
            }
        # enc_info 값이 있을 경우 결제 진행(결제인증 후 Ret_URL처리)
        else :
            post_data = {
                'req_tx' : f_get_parm(request.form['req_tx']), # 요청 종류         
                'res_cd' : f_get_parm(request.form['res_cd']), # 응답 코드
                'site_cd' : f_get_parm(request.form['site_cd']), # 사이트코드       
                'tran_cd' : f_get_parm(request.form['tran_cd']), # 트랜잭션 코드     
                'ordr_idxx' : f_get_parm(request.form['ordr_idxx']), # 쇼핑몰 주문번호   
                'good_name' : f_get_parm(request.form['good_name']), # 상품명            
                'good_mny' : f_get_parm(request.form['good_mny']), # 결제 금액       
                'buyr_name' : f_get_parm(request.form['buyr_name']), # 주문자명          
                'buyr_tel1' : f_get_parm(request.form['buyr_tel1']), # 주문자 전화번호   
                'buyr_tel2' : f_get_parm(request.form['buyr_tel2']), # 주문자 핸드폰 번호
                'buyr_mail' : f_get_parm(request.form['buyr_mail']), # 주문자 E-mail 주소
                'use_pay_method' : f_get_parm(request.form['use_pay_method']), # 결제 방법          
                'enc_info' : enc_info, # 암호화 정보       
                'enc_data' : f_get_parm(request.form['enc_data']), # 암호화 데이터     
                'param_opt_1' : '', # 기타 파라메터 추가 부분
                'param_opt_2' : '', # 기타 파라메터 추가 부분
                'param_opt_3' : ''  # 기타 파라메터 추가 부분
            }
            # 현금영수증 관련 데이터 처리(결제수단:계좌이체,포인트)
            if f_get_parm(request.form['use_pay_method']) == '010000000000' or f_get_parm(request.form['use_pay_method']) == '000100000000' :
                post_data.update({'cash_yn' : f_get_parm(request.form['cash_yn'])})
                # cash_yn(현금영수증발급여부) == 'Y' 인 경우
                if f_get_parm(request.form['cash_yn']) == 'Y':
                    post_data.update({'cash_tr_code' : f_get_parm(request.form['cash_tr_code'])})
    else :
        post_data = {
            'res_cd' : f_get_parm(request.form['res_cd']), # 응답 코드         
            'res_msg' : f_get_parm(request.form['res_msg']) # 응답 메세지
        }

    return render_template('mobile_sample/order_mobile.html', post_data=post_data)

# 결제요청 API (자동취소 로직 추가)
@app.route('/kcp_api_pay', methods=['POST'])
def kcp_api_pay():
    # 결과처리 POST DATA
    cash_yn = f_get_parm(request.form['cash_yn']) # 현금 영수증 등록 여부
    cash_tr_code = f_get_parm(request.form['cash_tr_code']) # 현금 영수증 발행 구분
    use_pay_method = f_get_parm(request.form['use_pay_method']) # 사용 결제수단
    ordr_idxx = f_get_parm(request.form['ordr_idxx']) # 주문번호
    
    post_data = {
        'cash_yn' : cash_yn,
        'cash_tr_code' : cash_tr_code,
        'use_pay_method' : use_pay_method,
        'ordr_idxx' : ordr_idxx
    }

    # 결제요청 API
    target_URL = 'https://stg-spl.kcp.co.kr/gw/enc/v1/payment' # 결제요청 API 개발환경
    #target_URL = 'https://spl.kcp.co.kr/gw/enc/v1/payment' # 결제요청 API 운영환경
    headers = {'Content-Type': 'application/json', 'charset': 'UTF-8'}

    site_cd = f_get_parm(request.form['site_cd'])

    # 결제 REQ DATA
    req_data = {
        'tran_cd' : f_get_parm(request.form['tran_cd']),
        'site_cd' : site_cd,
        'kcp_cert_info' : KCP_CERT_INFO,
        'enc_data' : f_get_parm(request.form['enc_data']),
        'enc_info' : f_get_parm(request.form['enc_info']),
        'ordr_mony' : '1' # ** 1 원은 실제로 업체에서 결제하셔야 될 원 금액을 넣어주셔야 합니다. 결제금액 유효성 검증 **
    }

    res = requests.post(target_URL, headers=headers, data=json.dumps(req_data, ensure_ascii=False, indent="\t").encode('utf8'))
     
    # ==========================================================================
    #      승인 결과 DB 처리 실패시 : 자동취소
    # --------------------------------------------------------------------------
    #      승인 결과를 DB 작업 하는 과정에서 정상적으로 승인된 건에 대해
    # DB 작업을 실패하여 DB update 가 완료되지 않은 경우, 자동으로
    #      승인 취소 요청을 하는 프로세스가 구성되어 있습니다.

    # DB 작업이 실패 한 경우, bSucc 라는 변수(String)의 값을 "false"
    #      로 설정해 주시기 바랍니다. (DB 작업 성공의 경우에는 "false" 이외의
    #      값을 설정하시면 됩니다.)
    # --------------------------------------------------------------------------
    bSucc = ''  # DB 작업 실패 또는 금액 불일치의 경우 "false" 로 세팅
    
    if bSucc == 'false' :

        res_data = json.loads(res.text)

        # 취소요청 API
        target_URL = 'https://stg-spl.kcp.co.kr/gw/mod/v1/cancel' # 취소요청 API 개발환경
        # target_URL = 'https://spl.kcp.co.kr/gw/mod/v1/cancel' # 취소요청 API 운영환경
        tno = res_data['tno']
        mod_type = 'STSC' # 전체취소
        
        # 서명대상데이터 (생성규칙)
        # 결제취소(Cancel) : site_cd^tno^mod_type
        cancel_sign_data = site_cd + '^' + tno + '^' + mod_type
        # 취소 서명데이터 생성
        kcp_sign_data = make_sign_data(cancel_sign_data)

        req_data = {
            'site_cd' : site_cd,
            'tno' : tno,
            'kcp_cert_info' : KCP_CERT_INFO,
            'kcp_sign_data' : kcp_sign_data,
            'mod_type' : mod_type,
            'mod_desc' : '가맹점 DB 처리 실패(자동취소)'
        }
        res = requests.post(target_URL, headers=headers, data=json.dumps(req_data, ensure_ascii=False, indent="\t").encode('utf8'))

    return render_template('kcp_api_pay.html', res_data=json.loads(res.text), req_data=req_data, post_data=post_data, bSucc=bSucc)

# null 값을 처리
def f_get_parm(val) :
    if val == None : val = ''
    return val

# 서명데이터 생성 예제
def make_sign_data(orgData) :
    # 개인키 READ
    # "splPrikeyPKCS8.pem" 은 테스트용 개인키
    key_file = open('C:\...\python_kcp_api_pay_sample\certificate\splPrikeyPKCS8.pem', 'r')
    key = key_file.read()
    key_file.close()

    # "changeit" 은 테스트용 개인키비밀번호
    password = 'changeit'.encode('utf-8')
    pkey = crypto.load_privatekey(crypto.FILETYPE_PEM, key, password)

    # 서명데이터생성
    sign = OpenSSL.crypto.sign(pkey, orgData, 'sha256')
    kcp_sign_data = base64.b64encode(sign).decode()
    
    return kcp_sign_data

if __name__ == '__main__':
    app.run(debug=True)
    