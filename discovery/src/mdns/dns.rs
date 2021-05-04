/// Opcode of DNS header
/// that generally should be 4 bit integer, 
/// but Rust does not have data type like that
#[derive(Debug, Clone)]
#[repr(u8)]
pub enum Opcode {
    Query = 0b0000,
    InverseQuery = 0x1,
    Status = 0x2,

    Notify = 0x4,
    Update = 0x5,
    DNSStatefulOperations = 0x6,


}

#[derive(Debug, Clone)]
pub struct Flags {
    /// Query/Response Flag
    /// The big is set to:
    /// 1 - Message is Response
    /// 0 - Message is Query
    pub qr: bool,

    /// Authoritative Answer Flag
    /// Meaningfull only in responses
    /// This bit is set to: 
    /// 1 - Response is authoritative
    /// 0 - Response is not authoritative
    pub aa: bool,

    /// Truncation Flag
    /// The bit is set to:
    /// 1 - Message was truncated due to its length being longer than max length in transport protocol
    /// 0 - Message was not truncated
    pub tc: bool,


    /// Recursion Desired
    /// This bit is set to:
    /// 1 - Recursion is desired
    /// 0 - Recursion is not desired
    pub rd: bool,
    
    /// Recursion Available
    /// This bit is set to:
    /// 1 - Recursion is available
    /// 0 - Recursion is not available
    pub ra: bool,
}

#[derive(Clone, Debug)]
#[repr(u8)]
/// Response of DNS Header
/// that generally should be 4 bit integer, 
/// but Rust does not have data type like that
pub enum ResponseCode {
    /// No error condition
    NoError,

    /// The name server was unable to interpret the query
    FormatError,

    /// The name server was unable to process this query due to a problem withthe name server
    ServerFailure,

    /// Meaningful only for responses from an authoritative name server
    /// this codesignifies that the domain name referenced in the query does not exist.
    NameError,

    /// The name server does not support the requested kind of query.
    NotImplemented,

    /// The name server refuses to perform the specified operation for policy reasons.
    Refused,
}



#[derive(Clone, Debug)]
pub struct Question {
    /// QNAME
    ///
    /// A domain name represented as a sequence of labels, where each label consists of a lengthoctet followed by that number of octets. 
    /// The domain name terminates with the zero lengthoctet for the null label of the root.
    pub name: String,
    
    /// QTYPE
    /// 
    /// A two octet code which specifies the type of the query.
    pub r#type: u16, // TODO: restrict this type

    /// QCLASS
    ///
    /// A two octet code that specifies the class of the query.
    pub class: u16,
}

#[derive(Clone, Debug)]
pub struct Record {
    /// NAME
    ///
    /// A domain name to which this resource record pertains.
    /// Null terminated string
    pub name: String,
    
    /// TYPE
    ///
    /// two octets containing one of the RR type codes.  
    /// This field specifies the meaning of the data in the RDATA field.
    // TODO: restrict this type
    pub r#type: u16, 

    /// CLASS
    ///
    /// A two octet code that specifies the class of the query.
    pub class: u16,

    /// TTL
    ///
    /// The number of seconds the results can be cached.
    pub ttl: u16,

    /// RDATA
    ///
    /// A variable length string of octets that describes the resource.
    /// The format of this information varies according to the TYPE and CLASS of the resource record.
    /// For example, the if the TYPE is A and the CLASS is IN, the RDATA field is a 4 octet ARPA Internet address.
    pub data: String,
}

#[derive(Clone, Debug)]
pub struct Header {
    /// A 16 bit identifier assigned by the program that generates any kind of query.
    /// This identifieris copied the corresponding reply and can be used by the requester to match up replies tooutstanding queries.
    pub id: u16,

    /// A 4 bit field that specifies opcode of packet
    pub opcode: Opcode,

    /// A 5 bit flags
    pub flags: Flags,

    /// Response code - this 4 bit field is set as part of responses. 
    pub response_code: ResponseCode,

    /// Number of entries in question section
    pub question_count: u16,

    /// Number of entries in answer section
    pub answer_count: u16,

    /// Number of entries in nameserver section
    pub nameserver_count: u16,

    /// Number of entries in records section
    pub records_count: u16,
}
