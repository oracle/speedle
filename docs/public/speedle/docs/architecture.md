#Architecture

Speedle is an authorization engine. It basically works like below diagram shows. 

<img src="../img/spdlarch.jpg" />


(1) users create/manage policies through PMS (Policy Management Service) API   
(2) PMS persists the policies in a policy repository, the repository could be a file or a database (TBD)   
(3) The policies are provisioned from policy repository to ADS (Authorization Decision Service, also known as ARS, Authorization Runtime Service)   
(4) Users systems invoke ADS API for authorization check   

If you are familar with XACML model, PMS serves as PAP, and ADS serves as PDP here.    



