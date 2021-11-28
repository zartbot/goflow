import pandas as pd
import re
cisco = pd.read_excel('./nf.xlsx',sheet_name="Sheet2")
cisco = cisco[(cisco['Val'].notnull()) & (cisco['Field Type'].notnull())]
regexp = re.compile(r'\(.*\)')
cisco['FullID']=cisco['Val'].replace(regexp,'').astype(int)
cisco = cisco[cisco['FullID']<65000]
cisco['ElementID'] = cisco['FullID'] - 32768
cisco['EnterpriseNo'] = 9
cisco['FieldType'] = cisco['Field Type'].apply(lambda x: x.replace('+','').replace(' ','_').replace('(','').replace(')','_').replace('%',''))
cisco = cisco[['FullID','ElementID','EnterpriseNo','FieldType','Len','Description']].reset_index(drop=True)
cisco['Description'].fillna('',inplace=True)
cisco['Len'] = cisco['Len'].apply(lambda x: str(x).replace('\n','none'))
#cisco['DataType'] = 'octetArray'
cisco.loc[cisco['FieldType'].str.contains('IPv4Address',regex=False),'DataType'] = 'ipv4Address'
cisco.loc[cisco['FieldType'].str.contains('IPv6Address',regex=False),'DataType'] = 'ipv6Address'
cisco.loc[cisco['FieldType'].str.contains('AddressIPv4',regex=False),'DataType'] = 'ipv4Address'
cisco.loc[cisco['FieldType'].str.contains('AddressIPv6',regex=False),'DataType'] = 'ipv6Address'
cisco.loc[(cisco['Description'].str.contains('addr',case=False) & cisco['Description'].str.contains('v4',case=False)) ,'DataType'] = 'ipv4Address'
cisco.loc[(cisco['Description'].str.contains('addr',case=False) & cisco['Description'].str.contains('v6',case=False)) ,'DataType'] = 'ipv6Address'

cisco.loc[cisco['FieldType'].str.contains('uri',case=False,regex=False),'DataType'] = 'string'
cisco.loc[cisco['FieldType'].str.contains('string',case=False,regex=False),'DataType'] = 'string'
cisco.loc[cisco['Description'].str.contains('string',case=False),'DataType'] = 'string'
cisco.loc[cisco['Len'].str.contains('string',case=False),'DataType'] = 'string'
cisco.loc[cisco['Len'].str.contains('Bitstring',case=False),'DataType'] = 'octetArray'
cisco.loc[cisco['Len'].str.contains('u8',case=False),'DataType'] = 'unsigned8'
cisco.loc[cisco['Len'].str.contains('u16',case=False),'DataType'] = 'unsigned16'
cisco.loc[cisco['Len'].str.contains('u32',case=False),'DataType'] = 'unsigned32'
cisco.loc[cisco['Len'].str.contains('u64',case=False),'DataType'] = 'unsigned64'
cisco.loc[cisco['FieldType'].str.contains('Port'),'DataType'] ='unsigned16'
cisco.loc[cisco['FieldType'].str.contains('Port',case=False) & (cisco['Len'] == '2'),'DataType'] = 'unsigned16'
cisco.loc[ (cisco['FieldType'].str.contains('count',case=False) &  (cisco['Len'] == '1')),'DataType'] = 'unsigned8'
cisco.loc[ (cisco['FieldType'].str.contains('count',case=False) &  (cisco['Len'] == '2')),'DataType'] = 'unsigned16'
cisco.loc[ (cisco['FieldType'].str.contains('count',case=False) &  (cisco['Len'] == '4')),'DataType'] = 'unsigned32'
cisco.loc[ (cisco['FieldType'].str.contains('count',case=False) &  (cisco['Len'] == '8')),'DataType'] = 'unsigned64'
cisco.loc[(cisco['Description'].str.contains('number',case=False) &  (cisco['Len'] == '1')),'DataType'] = 'unsigned8'
cisco.loc[(cisco['Description'].str.contains('number',case=False) &  (cisco['Len'] == '2')),'DataType'] = 'unsigned16'
cisco.loc[(cisco['Description'].str.contains('number',case=False) &  (cisco['Len'] == '4')),'DataType'] = 'unsigned32'
cisco.loc[(cisco['Description'].str.contains('number',case=False) &  (cisco['Len'] == '8')),'DataType'] = 'unsigned64'
cisco.loc[(cisco['Description'].str.contains('number',case=False) &  (cisco['Len'] == '4 or 8')),'DataType'] = 'unsigned64'
cisco.loc[(cisco['Description'].str.contains('msec',case=False)),'DataType'] = 'dateTimeMilliseconds'
cisco.loc[(cisco['Description'].str.contains('id',case=False)&  (cisco['Len'] == '1')),'DataType'] = 'unsigned8'
cisco.loc[(cisco['Description'].str.contains('id',case=False)&  (cisco['Len'] == '2')),'DataType'] = 'unsigned16'
cisco.loc[(cisco['Description'].str.contains('id',case=False)&  (cisco['Len'] == '4')),'DataType'] = 'unsigned32'
cisco.loc[(cisco['Description'].str.contains('id',case=False)&  (cisco['Len'] == '8')),'DataType'] = 'unsigned64'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == '1')),'DataType'] = 'unsigned8'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == '2')),'DataType'] = 'unsigned16'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == '4')),'DataType'] = 'unsigned32'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == '8')),'DataType'] = 'unsigned64'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == '32bits')),'DataType'] = 'unsigned32'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == '64bits')),'DataType'] = 'unsigned64'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == 'IPV4')),'DataType'] = 'ipv4Address'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == 'IPv4')),'DataType'] = 'ipv4Address'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == 'IPV6')),'DataType'] = 'ipv6Address'
cisco.loc[(cisco['DataType'].isna()&(cisco['Len'] == 'IPv6')),'DataType'] = 'ipv6Address'
cisco.loc[(cisco['DataType'].isna()),'DataType'] = 'octetArray'
cisco = cisco[['FullID','ElementID','EnterpriseNo','FieldType','DataType','Len','Description']]
cisco.to_csv('cisco_ie.csv')
iana_ie = pd.read_excel('./nf.xlsx',sheet_name="Sheet1")
iana_ie = iana_ie[(iana_ie['Val'].notnull()) & (iana_ie['Field Type'].notnull())]
regexp = re.compile(r'<.*')
iana_ie['FullID']=iana_ie['Val'].replace(regexp,'').astype(int)
iana_ie = iana_ie[iana_ie['FullID']<512]
iana_ie['ElementID'] = iana_ie['FullID'] 
iana_ie['EnterpriseNo'] = 0
iana_ie['FieldType'] = iana_ie['Field Type'].apply(lambda x: x.replace('+','').replace(' ','_').replace('(','').replace(')','_').replace('%','').replace('|',''))
iana_ie['Len'] = iana_ie['Length']
iana_ie = iana_ie[['FullID','ElementID','EnterpriseNo','FieldType','Len','Description']].reset_index(drop=True)
iana_ie_official = pd.read_csv('./ipfix-information-elements.csv')
iana_ie_official['DataType'] = iana_ie_official['Abstract Data Type'] 
dtype = iana_ie_official[['ElementID','DataType']]
iana_result = iana_ie.join(dtype,on="ElementID",rsuffix='r')[['FullID','ElementID','EnterpriseNo','FieldType','DataType','Len','Description']]
iana_result.to_csv('iana_ie.csv')
iana_ie_official = pd.read_csv('./ipfix-information-elements.csv')
iana_ie_official['DataType'] = iana_ie_official['Abstract Data Type'] 
iana_ie_official['FieldType']= iana_ie_official['Name']
iana_ie_official['FullID'] =  iana_ie_official['ElementID']
iana_ie_official['EnterpriseNo'] = 0
iana_ie_official['Len'] =  iana_ie_official['Units']
iana_ie_official[['FullID','ElementID','EnterpriseNo','FieldType','DataType','Len','Description']].to_csv('iana_ie2.csv')

