package ews

import (
	"encoding/xml"
	"errors"
)

type FindFoldersRequest struct {
	XMLName             struct{}            `xml:"m:FindFolder"`
	Traversal           string              `xml:"m:traversal"`
	FolderShape         *FolderShape        `xml:"m:FolderShape,omitempty"`
	IndexedPageItemView IndexedPageItemView `xml:"m:IndexedPageItemView"`
	ParentFolderIds     ParentFolderIds     `xml:"m:ParentFolderIds"`
}

type FolderShape struct {
	BaseShape            BaseShape            `xml:"t:BaseShape,omitempty"`
	AdditionalProperties AdditionalProperties `xml:"t:AdditionalProperties,omitempty"`
}

type ParentFolderIds struct {
	DistinguishedFolderId DistinguishedFolderId `xml:"t:DistinguishedFolderId"`
}

type findFolderResponseEnvelope struct {
	XMLName struct{}               `xml:"Envelope"`
	Body    findFolderResponseBody `xml:"Body"`
}

type findFolderResponseBody struct {
	FindFolderResponse FindFolderResponse `xml:"FindFolderResponse"`
}

type FindFolderResponse struct {
	ResponseMessage FindFolderResponseMessages `xml:"ResponseMessages"`
}

type FindFolderResponseMessages struct {
	FindFolderResponseMessage FindFolderResponseMessage `xml:"FindFolderResponseMessage"`
}

type FindFolderResponseMessage struct {
	Response
	RootFolder RootFolder `xml:"RootFolder"`
}

type RootFolder struct {
	TotalItemsInView        int     `xml:"TotalItemsInView,attr"`
	IncludesLastItemInRange bool    `xml:"IncludesLastItemInRange,attr"`
	Folders                 Folders `xml:"Folders"`
}

type Folders struct {
	Folder []Folder `xml:"Folder"`
}

type Folder struct {
	FolderId         FolderId                  `xml:"FolderId"`
	ParentFolderId   FindFoldersParentFolderId `xml:"ParentFolderId"`
	DisplayName      string                    `xml:"DisplayName"`
	TotalCount       int                       `xml:"TotalCount"`
	ChildFolderCount int                       `xml:"ChildFolderCount"`
	UnreadCount      int                       `xml:"UnreadCOunt"`
}

type FolderId struct {
	Id string `xml:"Id,attr"`
}

type FindFoldersParentFolderId struct {
	Id string `xml:"Id,attr"`
}

// FindFolders find mail foldes
func FindFolders(c Client, parentFolderName string) ([]Folder, error) {
	fieldURIs := []FieldURI{
		{FieldURI: "folder:ParentFolderId"},
		{FieldURI: "folder:FolderId"},
		{FieldURI: "folder:DisplayName"},
		{FieldURI: "folder:UnreadCount"},
		{FieldURI: "folder:TotalCount"},
	}
	req := &FindFoldersRequest{
		Traversal: "DEEP",
		FolderShape: &FolderShape{
			BaseShape: BaseShapeDefault,
			AdditionalProperties: AdditionalProperties{
				FieldURI: fieldURIs,
			},
		},
		IndexedPageItemView: IndexedPageItemView{
			MaxEntriesReturned: 1000,
			Offset:             0,
			BasePoint:          BasePointBeginning,
		},
		ParentFolderIds: ParentFolderIds{
			DistinguishedFolderId: DistinguishedFolderId{Id: "msgfolderroot"}},
	}

	resp, err := SendFindFolders(c, req)

	if err != nil {
		return nil, err
	}

	return resp.ResponseMessage.FindFolderResponseMessage.RootFolder.Folders.Folder, nil
}

func SendFindFolders(c Client, r *FindFoldersRequest) (*FindFolderResponse, error) {
	xmlBytes, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}

	bb, err := c.SendAndReceive(xmlBytes)
	if err != nil {
		return nil, err
	}

	var soapResp findFolderResponseEnvelope
	err = xml.Unmarshal(bb, &soapResp)
	if err != nil {
		return nil, err
	}

	if soapResp.Body.FindFolderResponse.ResponseMessage.FindFolderResponseMessage.ResponseClass == ResponseClassError {
		return nil, errors.New(soapResp.Body.FindFolderResponse.ResponseMessage.FindFolderResponseMessage.MessageText)
	}

	return &soapResp.Body.FindFolderResponse, nil
}
