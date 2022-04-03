### Functional Details

To add/modify one or more custom create attribute function, have to go config go file.

Add a new function by country in the map

 ```go
    var mapAttributes = map[string]func(data *Attributes) error{
    "AU": createAUAttributes,
    "BE": createBEAttributes,
    "CA": createCAAttributes,
    "EE": createEEAttributes,
    "GB": createGBAttributes,
    "FR": createFRAttributes,
    "-":  func(attributes *Attributes) error {
            return nil
        },
    }
 ```

And implement the function, for example createGBAttributes Function, it runs only United Kingdom country
 ```go
    func createGBAttributes(atrb *Attributes) error {
        err := validation.Errors{
            "bank_id":      validation.Validate(&atrb.BankID, validation.Required),
            "bic":          validation.Validate(&atrb.Bic, validation.Required),
            "bank_id_code": validation.Validate(&atrb.BankIDCode, validation.Required),
        }.Filter()
        if err != nil {
            return err
        }
        
        if atrb.BankIDCode != gbBankID {
            return errors.New("invalid Bank ID Code")
        }
        
        atrb.BaseCurrency = "GBP"
    
    return nil
    }
 ```