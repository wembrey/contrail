package db

import (
	"database/sql"
	"encoding/json"

	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/generated/models"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

const insertServiceInstanceQuery = "insert into `service_instance` (`uuid`,`virtual_router_id`,`max_instances`,`auto_scale`,`right_virtual_network`,`right_ip_address`,`management_virtual_network`,`left_virtual_network`,`left_ip_address`,`interface_list`,`ha_mode`,`availability_zone`,`auto_policy`,`key_value_pair`,`share`,`owner_access`,`owner`,`global_access`,`parent_uuid`,`parent_type`,`user_visible`,`permissions_owner_access`,`permissions_owner`,`other_access`,`group_access`,`group`,`last_modified`,`enable`,`description`,`creator`,`created`,`fq_name`,`display_name`,`annotations_key_value_pair`) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
const updateServiceInstanceQuery = "update `service_instance` set `uuid` = ?,`virtual_router_id` = ?,`max_instances` = ?,`auto_scale` = ?,`right_virtual_network` = ?,`right_ip_address` = ?,`management_virtual_network` = ?,`left_virtual_network` = ?,`left_ip_address` = ?,`interface_list` = ?,`ha_mode` = ?,`availability_zone` = ?,`auto_policy` = ?,`key_value_pair` = ?,`share` = ?,`owner_access` = ?,`owner` = ?,`global_access` = ?,`parent_uuid` = ?,`parent_type` = ?,`user_visible` = ?,`permissions_owner_access` = ?,`permissions_owner` = ?,`other_access` = ?,`group_access` = ?,`group` = ?,`last_modified` = ?,`enable` = ?,`description` = ?,`creator` = ?,`created` = ?,`fq_name` = ?,`display_name` = ?,`annotations_key_value_pair` = ?;"
const deleteServiceInstanceQuery = "delete from `service_instance` where uuid = ?"

// ServiceInstanceFields is db columns for ServiceInstance
var ServiceInstanceFields = []string{
	"uuid",
	"virtual_router_id",
	"max_instances",
	"auto_scale",
	"right_virtual_network",
	"right_ip_address",
	"management_virtual_network",
	"left_virtual_network",
	"left_ip_address",
	"interface_list",
	"ha_mode",
	"availability_zone",
	"auto_policy",
	"key_value_pair",
	"share",
	"owner_access",
	"owner",
	"global_access",
	"parent_uuid",
	"parent_type",
	"user_visible",
	"permissions_owner_access",
	"permissions_owner",
	"other_access",
	"group_access",
	"group",
	"last_modified",
	"enable",
	"description",
	"creator",
	"created",
	"fq_name",
	"display_name",
	"annotations_key_value_pair",
}

// ServiceInstanceRefFields is db reference fields for ServiceInstance
var ServiceInstanceRefFields = map[string][]string{

	"service_template": {
	// <common.Schema Value>

	},

	"instance_ip": {
		// <common.Schema Value>
		"interface_type",
	},
}

// ServiceInstanceBackRefFields is db back reference fields for ServiceInstance
var ServiceInstanceBackRefFields = map[string][]string{

	"port_tuple": {
		"uuid",
		"share",
		"owner_access",
		"owner",
		"global_access",
		"parent_uuid",
		"parent_type",
		"user_visible",
		"permissions_owner_access",
		"permissions_owner",
		"other_access",
		"group_access",
		"group",
		"last_modified",
		"enable",
		"description",
		"creator",
		"created",
		"fq_name",
		"display_name",
		"key_value_pair",
	},
}

const insertServiceInstanceServiceTemplateQuery = "insert into `ref_service_instance_service_template` (`from`, `to` ) values (?, ?);"

const insertServiceInstanceInstanceIPQuery = "insert into `ref_service_instance_instance_ip` (`from`, `to` ,`interface_type`) values (?, ?,?);"

// CreateServiceInstance inserts ServiceInstance to DB
func CreateServiceInstance(tx *sql.Tx, model *models.ServiceInstance) error {
	// Prepare statement for inserting data
	stmt, err := tx.Prepare(insertServiceInstanceQuery)
	if err != nil {
		return errors.Wrap(err, "preparing create statement failed")
	}
	defer stmt.Close()
	log.WithFields(log.Fields{
		"model": model,
		"query": insertServiceInstanceQuery,
	}).Debug("create query")
	_, err = stmt.Exec(string(model.UUID),
		string(model.ServiceInstanceProperties.VirtualRouterID),
		int(model.ServiceInstanceProperties.ScaleOut.MaxInstances),
		bool(model.ServiceInstanceProperties.ScaleOut.AutoScale),
		string(model.ServiceInstanceProperties.RightVirtualNetwork),
		string(model.ServiceInstanceProperties.RightIPAddress),
		string(model.ServiceInstanceProperties.ManagementVirtualNetwork),
		string(model.ServiceInstanceProperties.LeftVirtualNetwork),
		string(model.ServiceInstanceProperties.LeftIPAddress),
		common.MustJSON(model.ServiceInstanceProperties.InterfaceList),
		string(model.ServiceInstanceProperties.HaMode),
		string(model.ServiceInstanceProperties.AvailabilityZone),
		bool(model.ServiceInstanceProperties.AutoPolicy),
		common.MustJSON(model.ServiceInstanceBindings.KeyValuePair),
		common.MustJSON(model.Perms2.Share),
		int(model.Perms2.OwnerAccess),
		string(model.Perms2.Owner),
		int(model.Perms2.GlobalAccess),
		string(model.ParentUUID),
		string(model.ParentType),
		bool(model.IDPerms.UserVisible),
		int(model.IDPerms.Permissions.OwnerAccess),
		string(model.IDPerms.Permissions.Owner),
		int(model.IDPerms.Permissions.OtherAccess),
		int(model.IDPerms.Permissions.GroupAccess),
		string(model.IDPerms.Permissions.Group),
		string(model.IDPerms.LastModified),
		bool(model.IDPerms.Enable),
		string(model.IDPerms.Description),
		string(model.IDPerms.Creator),
		string(model.IDPerms.Created),
		common.MustJSON(model.FQName),
		string(model.DisplayName),
		common.MustJSON(model.Annotations.KeyValuePair))
	if err != nil {
		return errors.Wrap(err, "create failed")
	}

	stmtServiceTemplateRef, err := tx.Prepare(insertServiceInstanceServiceTemplateQuery)
	if err != nil {
		return errors.Wrap(err, "preparing ServiceTemplateRefs create statement failed")
	}
	defer stmtServiceTemplateRef.Close()
	for _, ref := range model.ServiceTemplateRefs {

		_, err = stmtServiceTemplateRef.Exec(model.UUID, ref.UUID)
		if err != nil {
			return errors.Wrap(err, "ServiceTemplateRefs create failed")
		}
	}

	stmtInstanceIPRef, err := tx.Prepare(insertServiceInstanceInstanceIPQuery)
	if err != nil {
		return errors.Wrap(err, "preparing InstanceIPRefs create statement failed")
	}
	defer stmtInstanceIPRef.Close()
	for _, ref := range model.InstanceIPRefs {

		if ref.Attr == nil {
			ref.Attr = models.MakeServiceInterfaceTag()
		}

		_, err = stmtInstanceIPRef.Exec(model.UUID, ref.UUID, string(ref.Attr.InterfaceType))
		if err != nil {
			return errors.Wrap(err, "InstanceIPRefs create failed")
		}
	}

	log.WithFields(log.Fields{
		"model": model,
	}).Debug("created")
	return err
}

func scanServiceInstance(values map[string]interface{}) (*models.ServiceInstance, error) {
	m := models.MakeServiceInstance()

	if value, ok := values["uuid"]; ok {

		castedValue := common.InterfaceToString(value)

		m.UUID = castedValue

	}

	if value, ok := values["virtual_router_id"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ServiceInstanceProperties.VirtualRouterID = castedValue

	}

	if value, ok := values["max_instances"]; ok {

		castedValue := common.InterfaceToInt(value)

		m.ServiceInstanceProperties.ScaleOut.MaxInstances = castedValue

	}

	if value, ok := values["auto_scale"]; ok {

		castedValue := common.InterfaceToBool(value)

		m.ServiceInstanceProperties.ScaleOut.AutoScale = castedValue

	}

	if value, ok := values["right_virtual_network"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ServiceInstanceProperties.RightVirtualNetwork = castedValue

	}

	if value, ok := values["right_ip_address"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ServiceInstanceProperties.RightIPAddress = models.IpAddressType(castedValue)

	}

	if value, ok := values["management_virtual_network"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ServiceInstanceProperties.ManagementVirtualNetwork = castedValue

	}

	if value, ok := values["left_virtual_network"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ServiceInstanceProperties.LeftVirtualNetwork = castedValue

	}

	if value, ok := values["left_ip_address"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ServiceInstanceProperties.LeftIPAddress = models.IpAddressType(castedValue)

	}

	if value, ok := values["interface_list"]; ok {

		json.Unmarshal(value.([]byte), &m.ServiceInstanceProperties.InterfaceList)

	}

	if value, ok := values["ha_mode"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ServiceInstanceProperties.HaMode = models.AddressMode(castedValue)

	}

	if value, ok := values["availability_zone"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ServiceInstanceProperties.AvailabilityZone = castedValue

	}

	if value, ok := values["auto_policy"]; ok {

		castedValue := common.InterfaceToBool(value)

		m.ServiceInstanceProperties.AutoPolicy = castedValue

	}

	if value, ok := values["key_value_pair"]; ok {

		json.Unmarshal(value.([]byte), &m.ServiceInstanceBindings.KeyValuePair)

	}

	if value, ok := values["share"]; ok {

		json.Unmarshal(value.([]byte), &m.Perms2.Share)

	}

	if value, ok := values["owner_access"]; ok {

		castedValue := common.InterfaceToInt(value)

		m.Perms2.OwnerAccess = models.AccessType(castedValue)

	}

	if value, ok := values["owner"]; ok {

		castedValue := common.InterfaceToString(value)

		m.Perms2.Owner = castedValue

	}

	if value, ok := values["global_access"]; ok {

		castedValue := common.InterfaceToInt(value)

		m.Perms2.GlobalAccess = models.AccessType(castedValue)

	}

	if value, ok := values["parent_uuid"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ParentUUID = castedValue

	}

	if value, ok := values["parent_type"]; ok {

		castedValue := common.InterfaceToString(value)

		m.ParentType = castedValue

	}

	if value, ok := values["user_visible"]; ok {

		castedValue := common.InterfaceToBool(value)

		m.IDPerms.UserVisible = castedValue

	}

	if value, ok := values["permissions_owner_access"]; ok {

		castedValue := common.InterfaceToInt(value)

		m.IDPerms.Permissions.OwnerAccess = models.AccessType(castedValue)

	}

	if value, ok := values["permissions_owner"]; ok {

		castedValue := common.InterfaceToString(value)

		m.IDPerms.Permissions.Owner = castedValue

	}

	if value, ok := values["other_access"]; ok {

		castedValue := common.InterfaceToInt(value)

		m.IDPerms.Permissions.OtherAccess = models.AccessType(castedValue)

	}

	if value, ok := values["group_access"]; ok {

		castedValue := common.InterfaceToInt(value)

		m.IDPerms.Permissions.GroupAccess = models.AccessType(castedValue)

	}

	if value, ok := values["group"]; ok {

		castedValue := common.InterfaceToString(value)

		m.IDPerms.Permissions.Group = castedValue

	}

	if value, ok := values["last_modified"]; ok {

		castedValue := common.InterfaceToString(value)

		m.IDPerms.LastModified = castedValue

	}

	if value, ok := values["enable"]; ok {

		castedValue := common.InterfaceToBool(value)

		m.IDPerms.Enable = castedValue

	}

	if value, ok := values["description"]; ok {

		castedValue := common.InterfaceToString(value)

		m.IDPerms.Description = castedValue

	}

	if value, ok := values["creator"]; ok {

		castedValue := common.InterfaceToString(value)

		m.IDPerms.Creator = castedValue

	}

	if value, ok := values["created"]; ok {

		castedValue := common.InterfaceToString(value)

		m.IDPerms.Created = castedValue

	}

	if value, ok := values["fq_name"]; ok {

		json.Unmarshal(value.([]byte), &m.FQName)

	}

	if value, ok := values["display_name"]; ok {

		castedValue := common.InterfaceToString(value)

		m.DisplayName = castedValue

	}

	if value, ok := values["annotations_key_value_pair"]; ok {

		json.Unmarshal(value.([]byte), &m.Annotations.KeyValuePair)

	}

	if value, ok := values["ref_service_template"]; ok {
		var references []interface{}
		stringValue := common.InterfaceToString(value)
		json.Unmarshal([]byte("["+stringValue+"]"), &references)
		for _, reference := range references {
			referenceMap, ok := reference.(map[string]interface{})
			if !ok {
				continue
			}
			if referenceMap["to"] == "" {
				continue
			}
			referenceModel := &models.ServiceInstanceServiceTemplateRef{}
			referenceModel.UUID = common.InterfaceToString(referenceMap["to"])
			m.ServiceTemplateRefs = append(m.ServiceTemplateRefs, referenceModel)

		}
	}

	if value, ok := values["ref_instance_ip"]; ok {
		var references []interface{}
		stringValue := common.InterfaceToString(value)
		json.Unmarshal([]byte("["+stringValue+"]"), &references)
		for _, reference := range references {
			referenceMap, ok := reference.(map[string]interface{})
			if !ok {
				continue
			}
			if referenceMap["to"] == "" {
				continue
			}
			referenceModel := &models.ServiceInstanceInstanceIPRef{}
			referenceModel.UUID = common.InterfaceToString(referenceMap["to"])
			m.InstanceIPRefs = append(m.InstanceIPRefs, referenceModel)

			attr := models.MakeServiceInterfaceTag()
			referenceModel.Attr = attr

		}
	}

	if value, ok := values["backref_port_tuple"]; ok {
		var childResources []interface{}
		stringValue := common.InterfaceToString(value)
		json.Unmarshal([]byte("["+stringValue+"]"), &childResources)
		for _, childResource := range childResources {
			childResourceMap, ok := childResource.(map[string]interface{})
			if !ok {
				continue
			}
			if childResourceMap["uuid"] == "" {
				continue
			}
			childModel := models.MakePortTuple()
			m.PortTuples = append(m.PortTuples, childModel)

			if propertyValue, ok := childResourceMap["uuid"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.UUID = castedValue

			}

			if propertyValue, ok := childResourceMap["share"]; ok && propertyValue != nil {

				json.Unmarshal(common.InterfaceToBytes(propertyValue), &childModel.Perms2.Share)

			}

			if propertyValue, ok := childResourceMap["owner_access"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToInt(propertyValue)

				childModel.Perms2.OwnerAccess = models.AccessType(castedValue)

			}

			if propertyValue, ok := childResourceMap["owner"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.Perms2.Owner = castedValue

			}

			if propertyValue, ok := childResourceMap["global_access"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToInt(propertyValue)

				childModel.Perms2.GlobalAccess = models.AccessType(castedValue)

			}

			if propertyValue, ok := childResourceMap["parent_uuid"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.ParentUUID = castedValue

			}

			if propertyValue, ok := childResourceMap["parent_type"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.ParentType = castedValue

			}

			if propertyValue, ok := childResourceMap["user_visible"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToBool(propertyValue)

				childModel.IDPerms.UserVisible = castedValue

			}

			if propertyValue, ok := childResourceMap["permissions_owner_access"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToInt(propertyValue)

				childModel.IDPerms.Permissions.OwnerAccess = models.AccessType(castedValue)

			}

			if propertyValue, ok := childResourceMap["permissions_owner"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.IDPerms.Permissions.Owner = castedValue

			}

			if propertyValue, ok := childResourceMap["other_access"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToInt(propertyValue)

				childModel.IDPerms.Permissions.OtherAccess = models.AccessType(castedValue)

			}

			if propertyValue, ok := childResourceMap["group_access"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToInt(propertyValue)

				childModel.IDPerms.Permissions.GroupAccess = models.AccessType(castedValue)

			}

			if propertyValue, ok := childResourceMap["group"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.IDPerms.Permissions.Group = castedValue

			}

			if propertyValue, ok := childResourceMap["last_modified"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.IDPerms.LastModified = castedValue

			}

			if propertyValue, ok := childResourceMap["enable"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToBool(propertyValue)

				childModel.IDPerms.Enable = castedValue

			}

			if propertyValue, ok := childResourceMap["description"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.IDPerms.Description = castedValue

			}

			if propertyValue, ok := childResourceMap["creator"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.IDPerms.Creator = castedValue

			}

			if propertyValue, ok := childResourceMap["created"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.IDPerms.Created = castedValue

			}

			if propertyValue, ok := childResourceMap["fq_name"]; ok && propertyValue != nil {

				json.Unmarshal(common.InterfaceToBytes(propertyValue), &childModel.FQName)

			}

			if propertyValue, ok := childResourceMap["display_name"]; ok && propertyValue != nil {

				castedValue := common.InterfaceToString(propertyValue)

				childModel.DisplayName = castedValue

			}

			if propertyValue, ok := childResourceMap["key_value_pair"]; ok && propertyValue != nil {

				json.Unmarshal(common.InterfaceToBytes(propertyValue), &childModel.Annotations.KeyValuePair)

			}

		}
	}

	return m, nil
}

// ListServiceInstance lists ServiceInstance with list spec.
func ListServiceInstance(tx *sql.Tx, spec *common.ListSpec) ([]*models.ServiceInstance, error) {
	var rows *sql.Rows
	var err error
	//TODO (check input)
	spec.Table = "service_instance"
	spec.Fields = ServiceInstanceFields
	spec.RefFields = ServiceInstanceRefFields
	spec.BackRefFields = ServiceInstanceBackRefFields
	result := models.MakeServiceInstanceSlice()
	query, columns, values := common.BuildListQuery(spec)
	log.WithFields(log.Fields{
		"listSpec": spec,
		"query":    query,
	}).Debug("select query")
	rows, err = tx.Query(query, values...)
	if err != nil {
		return nil, errors.Wrap(err, "select query failed")
	}
	defer rows.Close()
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "row error")
	}
	for rows.Next() {
		valuesMap := map[string]interface{}{}
		values := make([]interface{}, len(columns))
		valuesPointers := make([]interface{}, len(columns))
		for _, index := range columns {
			valuesPointers[index] = &values[index]
		}
		if err := rows.Scan(valuesPointers...); err != nil {
			return nil, errors.Wrap(err, "scan failed")
		}
		for column, index := range columns {
			val := valuesPointers[index].(*interface{})
			valuesMap[column] = *val
		}
		log.WithFields(log.Fields{
			"valuesMap": valuesMap,
		}).Debug("valueMap")
		m, err := scanServiceInstance(valuesMap)
		if err != nil {
			return nil, errors.Wrap(err, "scan row failed")
		}
		result = append(result, m)
	}
	return result, nil
}

// ShowServiceInstance shows ServiceInstance resource
func ShowServiceInstance(tx *sql.Tx, uuid string) (*models.ServiceInstance, error) {
	list, err := ListServiceInstance(tx, &common.ListSpec{
		Filter: map[string]interface{}{"uuid": uuid},
		Limit:  1})
	if len(list) == 0 {
		return nil, errors.Wrap(err, "show query failed")
	}
	return list[0], err
}

// UpdateServiceInstance updates a resource
func UpdateServiceInstance(tx *sql.Tx, uuid string, model *models.ServiceInstance) error {
	//TODO(nati) support update
	return nil
}

// DeleteServiceInstance deletes a resource
func DeleteServiceInstance(tx *sql.Tx, uuid string) error {
	stmt, err := tx.Prepare(deleteServiceInstanceQuery)
	if err != nil {
		return errors.Wrap(err, "preparing delete query failed")
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid)
	if err != nil {
		return errors.Wrap(err, "delete failed")
	}
	log.WithFields(log.Fields{
		"uuid": uuid,
	}).Debug("deleted")
	return nil
}