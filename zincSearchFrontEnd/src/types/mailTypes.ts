export interface Mail{
	ID: number
    'Message-ID': string
	Date: string 
	From: string 
	To: string 
	Subject: string
	Cc: string
	'Mime-Version': string
	'Content-Type': string 
	'Content-Transfer-Encoding': string 
	Bcc: string 
	'X-From': string 
	'X-To': string
	'X-cc': string 
	'X-bcc': string 
	'X-Folder': string 
	'X-Origin':string 
	'X-FileName': string
	Body: string 
}