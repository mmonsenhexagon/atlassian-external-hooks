package com.ngs.stash.externalhooks.hook;

import org.springframework.stereotype.Component;
import org.springframework.beans.factory.annotation.Autowired;

import com.atlassian.upm.api.license.PluginLicenseManager;
import com.atlassian.upm.api.license.entity.LicenseError;
import com.atlassian.upm.api.license.entity.PluginLicense;
import com.atlassian.upm.api.util.Option;

@Component
public class License {
	@Autowired
	private PluginLicenseManager pluginLicenseManager;

	public boolean isValid() {
		Option<PluginLicense> licenseOption = pluginLicenseManager.getLicense();
		if (!licenseOption.isDefined()) {
			return true;
		}

		PluginLicense pluginLicense = licenseOption.get();
		return pluginLicense.isValid();
	}

	//public LicenseStatus getStatus() {
	//    Option<PluginLicense> licenseOption = pluginLicenseManager.getLicense();
	//    LicenseStatus licenseStatus;
	//    if (licenseOption.isDefined()) {
	//        PluginLicense pluginLicense = licenseOption.get();
	//        if (pluginLicense.isValid()) {
	//            licenseStatus = LicenseStatus.OK;
	//        }
	//        Option<LicenseError> error = pluginLicense.getError();
	//        licenseStatus = LicenseStatus.valueOf(error.get().name());
	//    } else {
	//        licenseStatus = LicenseStatus.OK;
	//    }

	//    return licenseStatus;
	//}
}
