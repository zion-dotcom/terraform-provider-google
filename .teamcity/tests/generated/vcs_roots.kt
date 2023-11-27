/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is auto-generated with mmv1, any changes made here will be overwritten

package tests

import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleRootProject

class VcsTests {
    @Test
    fun buildsHaveCleanCheckOut() {
        val project = googleRootProject(testConfiguration())
        project.buildTypes.forEach { bt ->
            assertTrue("Build '${bt.id}' doesn't use clean checkout", bt.vcs.cleanCheckout)
        }
    }
}
